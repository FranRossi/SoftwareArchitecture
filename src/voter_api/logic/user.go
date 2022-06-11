package logic

//func RegisterUser(id string, username string, password string, role string) (*domain.User, error) {
//	user, err := createUser(id, username, password, role)
//	if err != nil {
//		log.Fatal(err)
//		return nil, err
//	}
//	return user, err
//}
//
//func createUser(id string, username string, password string, role string) (*domain.User, error) {
//	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	//if err != nil {
//	//	return nil, fmt.Errorf("cannot hash password: %w", err)
//	//}
//	//
//	//user := &domain.User{
//	//	Id:             id,
//	//	Username:       username,
//	//	HashedPassword: string(hashedPassword),
//	//	Role:           role,
//	//}
//	//return StoreUser(user)
//	return nil, nil
//}
//
//func StoreUser(user *domain.User) (*domain.User, error) {
//	err := repository.RegisterUser(user)
//	if err != nil {
//		return nil, fmt.Errorf("user cannot be created: %w", err)
//	}
//	return user, nil
//}
//
//func IsCorrectPassword(user *domain.User, password string) bool {
//	//err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
//	//return err == nil
//	return true
//}
