package funcs

import (
	"database/sql"
	"errors"
	"fmt"

	database "funcs/internal/database"
	pkg "funcs/pkg"
)

func InsertEvent(event Event) (int, error) {
	query := "INSERT INTO posts_event (user_id, group_id, title, description, eventDate) VALUES (?,?,?,?,?)"
	db := database.GetDb()

	result, err := db.Exec(query, event.CreatorId, event.GroupId, event.Title, event.Descreption, event.Time)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func InsertGroupInfo(group Group) error {
	db := database.GetDb()

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}

	query := "INSERT INTO groups (creator_id, title, description, imagePath) VALUES (?,?,?,?)"
	result, err := tx.Exec(query, group.CreatorId, group.Title, group.Descreption, group.Path)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	group.Id = int(id)

	query = "INSERT INTO group_members (group_id,user_id,role) VALUES (?,?,?)"
	_, err = tx.Exec(query, group.Id, group.CreatorId, StatusCreator)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func UnFollow(userId, groupId int) error {
	// error?
	query := "DELETE FROM groups_members WHERE groupId=? AND userId=?"
	db := database.GetDb()
	_, err := db.Exec(query, groupId, userId)
	return err
}

func IsFollowing(userId, groupId int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)"
	db := database.GetDb()

	var exists bool
	err := db.QueryRow(query, groupId, userId).Scan(&exists)
	if err != nil {
		// if err==sql.ErrNoRows{
		// 	err=nil
		// }
		fmt.Println("err ahabchi", err)
		return false, err
	}
	return exists, nil
}

func GetuserGroupDB(userId int) ([]Group, error) {
	var groups []Group
	query := "SELECT group_id FROM group_members WHERE user_id=?"
	db := database.GetDb()
	rows, err := db.Query(query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return groups, err
	}

	for rows.Next() {

		groupId := 0
		err = rows.Scan(&groupId)
		if err != nil {
			return groups, err
		}
		group, err := GetGroupInfo(groupId)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func GetGroupInfo(groupId int) (Group, error) {
	var group Group
	group.Id = groupId
	db := database.GetDb()
	query := "SELECT creator_id, title, createdAt, description,imagePath FROM groups WHERE id=?"
	err := db.QueryRow(query, group.Id).Scan(
		&group.CreatorId,
		&group.Title,
		&group.CreatedAt,
		&group.Descreption,
		&group.Path,
	)
	if err != nil {
		return group, err
	}

	group.NbrMembers, err = CountMembers(groupId)
	return group, err
}

func CheckGroupExists(groupId int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM groups WHERE id = ?)"
	db := database.GetDb()
	var exist bool = false
	err := db.QueryRow(query, groupId).Scan(&exist)
	return exist, err
}

func CountMembers(groupId int) (int, error) {
	var nbr int
	db := database.GetDb()
	query := "SELECT COUNT(*) FROM group_members WHERE group_id=?"
	err := db.QueryRow(query, groupId).Scan(&nbr)
	if err != nil && err != sql.ErrNoRows {
		return nbr, err
	}
	return nbr, nil
}

func GetuserInfo(id int) (SUser, error) {
	query := "SELECT id, firstName, lastName, nickname, avatar FROM users WHERE id=?"
	db := database.GetDb()
	var user SUser

	err := db.QueryRow(query, id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Nickname, &user.Path)
	if err == sql.ErrNoRows {
		return user, nil
	}

	return user, err
}

func CheckIfMember(userId, groupId int) (STATUS, error) {
	var status STATUS
	query := "SELECT role,joined_at FROM group_members WHERE user_id=? AND group_id=?"
	db := database.GetDb()

	err := db.QueryRow(query, userId, groupId).Scan(&status.Status, &status.Since)
	if errors.Is(err, sql.ErrNoRows) {
		err = pkg.ErrNotMember
	}
	return status, err
}

func GetGroupCreator(groupId int) (SUser, error) {
	db := database.GetDb()
	query := "SELECT creator_id FROM groups WHERE id=?"
	var user SUser

	err := db.QueryRow(query, groupId).Scan(&user.Id)
	if err != nil {
		return user, err
	}

	return GetuserInfo(user.Id)
}

func CheckIfInveted(userId, groupId int) (STATUS, error) {
	var status STATUS
	query := "SELECT invited_by,invited_at FROM invitations WHERE invited_user=? AND group_id=? AND status='pending'"
	db := database.GetDb()

	err := db.QueryRow(query, userId, groupId).Scan(&status.Sender.Id, &status.Since)
	if errors.Is(err, sql.ErrNoRows) {
		err = pkg.ErrNoInvitation
	}
	return status, err
}

func CheckIfSendRequest(userId, groupId int) (STATUS, error) {
	var status STATUS
	status.Status = StatusSendRequest
	query := "SELECT requested_at FROM join_requests WHERE user_id=? AND group_id=? AND status='pending'"
	db := database.GetDb()

	err := db.QueryRow(query, userId, groupId).Scan(&status.Since)
	if errors.Is(err, sql.ErrNoRows) {
		err = pkg.ErrNoRequest
	}
	return status, err
}

func InviteUser(senderId, reciverId, groupId int) error {
	db := database.GetDb()
	query := "INSERT OR REPLACE INTO invitations (group_id, invited_by, invited_user, status) VALUES (?, ?, ?, 'pending');"
	_, err := db.Exec(query, groupId, senderId, reciverId)
	return err
}

func SendRequest(userId, groupId int) error {
	db := database.GetDb()
	query := "INSERT OR REPLACE INTO join_requests (group_id, user_id, status) VALUES (?, ?, 'pending');"
	_, err := db.Exec(query, groupId, userId)
	return err
}

func Accept(userId, groupId int) error {
	db := database.GetDb()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err := AcceptAllInvitationAndRequest(tx, userId, groupId); err != nil {
		return fmt.Errorf("failed to accept invitations/requests: %w", err)
	}

	if err := InsertMember(tx, groupId, userId); err != nil {
		return fmt.Errorf("failed to insert member: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func AcceptAllInvitationAndRequest(tx *sql.Tx, userId, groupId int) error {
	query := "UPDATE invitations SET status = 'accepted' WHERE group_id = ? AND invited_user = ? AND status = 'pending';"
	if _, err := tx.Exec(query, groupId, userId); err != nil {
		return fmt.Errorf("failed to update invitations: %w", err)
	}

	query = "UPDATE join_requests SET status = 'accepted' WHERE group_id = ? AND user_id = ? AND status = 'pending';"
	if _, err := tx.Exec(query, groupId, userId); err != nil {
		return fmt.Errorf("failed to update join requests: %w", err)
	}

	return nil
}

func InsertMember(tx *sql.Tx, groupId, userId int) error {
	checkQuery := "SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?"
	var count int
	if err := tx.QueryRow(checkQuery, groupId, userId).Scan(&count); err != nil {
		return fmt.Errorf("failed to check membership: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("user %d is already a member of group %d", userId, groupId)
	}

	insertQuery := "INSERT INTO group_members(group_id, user_id) VALUES(?,?)"
	if _, err := tx.Exec(insertQuery, groupId, userId); err != nil {
		return fmt.Errorf("failed to insert into groups_members: %w", err)
	}

	return nil
}

func CheckIfCreator(userId, groupId int) (bool, error) {
	quety := `SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ? AND role = 'creator')`
	db := database.GetDb()
	var exist bool = false
	err := db.QueryRow(quety, groupId, userId).Scan(&exist)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exist, nil
}

func GetJoinedGroups(userid int) ([]SGroup, error) {
	query := `SELECT g.id, g.title, g.imagePath, COUNT(m2.user_id) as member_count 
    FROM groups g 
    JOIN group_members m ON g.id = m.group_id
    LEFT JOIN group_members m2 ON g.id = m2.group_id
    WHERE m.user_id = ?
    GROUP BY g.id, g.title, g.imagePath`

	var groups []SGroup
	db := database.GetDb()
	rows, err := db.Query(query, userid)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return groups, err
	}
	defer rows.Close()

	for rows.Next() {
		var g SGroup
		err = rows.Scan(&g.Id, &g.Title, &g.Path, &g.MemberCount)
		if err != nil {
			return groups, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func Reject(userId, groupId int) error {
	query := "DELETE  FROM invitations WHERE (group_id = ? AND invited_user = ? AND status = ?) ;"
	db := database.GetDb()
	_, err := db.Exec(query, groupId, userId, "pending")
	if err == sql.ErrNoRows {
		err = nil
	}

	query = "DELETE  FROM join_requests WHERE (group_id = ? AND user_id = ? AND status = ?) ;"
	_, err = db.Exec(query, groupId, userId, "pending")
	return err
}

func GetGroupEvents(groupId int) ([]Event, error) {
	query := "SELECT id, user_id, group_id, title, description, eventDate FROM posts_event WHERE group_id = ? ORDER BY eventDate DESC;"
	db := database.GetDb()
	rows, err := db.Query(query, groupId)
	if err != nil {
		return []Event{}, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.CreatorId, &event.GroupId, &event.Title, &event.Descreption, &event.Time)
		if err != nil {
			return []Event{}, err
		}

		event.User, err = GetuserInfo(event.CreatorId)
		if err != nil {
			return []Event{}, err
		}
		event.CountAttends, err = GetEventAttendeesCount(event.Id)
		if err != nil {
			return []Event{}, err
		}
		events = append(events, event)
	}

	return events, nil
}

func JoinEvent(member_id, groupId, eventId int) error {
	query := "INSERT INTO EventStatus (member_id, group_id, event_id) VALUES (?, ?, ?)"
	db := database.GetDb()
	_, err := db.Exec(query, member_id, groupId, eventId)
	return err
}

func DeleteAttending(member_id, groupId, eventId int) error {
	query := "DELETE FROM Eventstatus WHERE member_id = ? AND group_id = ? AND event_id = ?"
	db := database.GetDb()
	_, err := db.Exec(query, member_id, groupId, eventId)
	return err
}

func CheckEventRegistration(member_id, groupId, eventId int) (bool, error) {
	db := database.GetDb()
	var count int
	query := "SELECT COUNT(*) FROM EventStatus WHERE member_id = ? AND group_id = ? AND event_id = ?"
	err := db.QueryRow(query, member_id, groupId, eventId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetEventAttendeesCount(eventId int) (int, error) {
	db := database.GetDb()
	var count int
	query := "SELECT COUNT(*) FROM EventStatus WHERE event_id = ?"
	err := db.QueryRow(query, eventId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

