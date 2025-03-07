package funcs

func Follow(id, targetId int) error {
	profiletype, err := GetProfileType(targetId)
	if err != nil {
		return err
	}
	result, _ := GetNowRequest(id, targetId)
	if result {
		return nil
	}
	if profiletype && !result {
		return AddFollower(id, targetId, 0)
	} 
	return AddFollower(id, targetId, 1)
}

func UnFollow(id, targetId int) error {
	return DeleteFollower(id, targetId)
}

func GetProfileDB(targetId, userId int) (Profile, error) {
	profile, err := GetProfileInfo(targetId, userId)
	if err != nil {
		return profile, err
	}

	profile.NbrPosts, err = CountPosts(targetId)
	if err != nil {
		return profile, err
	}

	// profile.Followers, err = CountFollowers(targetId)
	// if err != nil {
	// 	return profile, err
	// }
	// profile.IsFollow, err = IsFriends2(targetId, userId)
	// if err != nil {
	// 	return profile, err
	// }

	return profile, err
}
