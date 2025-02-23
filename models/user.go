type User struct {
	ID       primitive.ObjectID `bson:"_id",omitempty`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

func CreateUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Println("Error hashing password", err)
		return err
	}
	user.Password = string(hashedPassword)
	_, err = config.DB.Collection("users").InsertOne(ctx, user)
	return err
}

func FindUserByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := config.DB.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}