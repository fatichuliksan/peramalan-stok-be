package usecase

import (
	"gorm.io/gorm"
)

// PingInterface ...
type PingInterface interface {
	Ping(text string) (string, error)
}

type ping struct {
	DB *gorm.DB
}

// NewPing ...
func NewPing(db *gorm.DB) PingInterface {
	return &ping{
		DB: db,
	}
}

// Ping ...
func (t *ping) Ping(text string) (string, error) {

	var err error
	// for i := 0; i < 1; i++ {
	// 	go func(db postgre.Interface, i int) {
	// 		conn, err := t.DB.Connection()
	// 		if err == nil {
	// 			userRepo := repository.NewUserRepo(conn)
	// 			user, err := userRepo.FindOne(model.User{
	// 				Email: "director-sales-kotamalang@nta.co.id",
	// 			})
	// 			logger.Default().Println(err)
	// 			logger.Default().Println(fmt.Sprintf("%v %v", i, user.Email))
	// 		}
	// 	}(t.DB, i)
	// }

	// conn, err := t.DB.Connection()
	// if err != nil {
	// 	return text, err
	// }
	// userRepo := repository.NewUser(conn)
	// user, err := userRepo.Create(model.User{
	// 	Name: "test",
	// })
	// if err != nil {
	// 	return text, err
	// }
	return text, err
}
