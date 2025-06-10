package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/epic55/BankAppNew/internal/models"
	"github.com/epic55/BankAppNew/internal/repository"

	pb "github.com/epic55/BankAppNew/buyingGRPC"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ServiceInterface interface {
	GetUser(id int64) (*models.User, error)
	//Buying(ctx context.Context, req *pb.BuyingRequest) (*pb.BuyingReply, error)
	pb.BuyingServer
}

type ServiceStruct struct {
	repo repository.RepositoryInterface
	pb.UnimplementedBuyingServer
}

func NewService(repo repository.RepositoryInterface) ServiceInterface {
	return &ServiceStruct{repo: repo}
}

func (s *ServiceStruct) GetUser(id int64) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *ServiceStruct) Buying(ctx context.Context, req *pb.BuyingRequest) (*pb.BuyingReply, error) {
	var m sync.Mutex

	date := time.Now()
	date1 := date.Format("2006-01-02 15:04:05")

	account, err := s.repo.BuyingRepo(int(req.UserId))
	if err != nil {
		log.Println("Error fetching account:", err)
		return nil, fmt.Errorf("error fetching account: %v", err)
	}

	if checkPin(int(req.Pin)) {
		if !account.Blocked {

			if account.Balance >= float64(req.Price) {
				updatedBalance := account.Balance - float64(req.Price)

				typeofoperation2 := "buying"
				m.Lock()
				s.repo.UpdateAccount(updatedBalance, float64(req.Price), int(req.UserId), account.Currency, typeofoperation2, date1)
				m.Unlock()

				//typeofoperation := "buying"
				//h.R.UpdateHistory(typeofoperation, account.Name, account.Currency, float64(req.Price), date1)

			} else {
				NotEnoughMoney()
				return nil, fmt.Errorf("Not Enough Money")
			}

		} else {
			AccountIsBlocked(account.Name, account.Id)
			return nil, fmt.Errorf("Account Is Blocked")
		}
	} else {
		fmt.Println("PIN is incorrect")
		return nil, fmt.Errorf("PIN is incorrect")
	}

	log.Printf("Received %v for buying", req.GetPrice())
	return &pb.BuyingReply{Message: "User paid " + strconv.Itoa(int(req.GetPrice())) + " tg for ticket"}, nil
}

func checkPin(item int) bool { //CHECK PIN
	const pin = 1234
	if pin == item {
		return true
	}
	return false
}

func NotEnoughMoney() *pb.BuyingReply {
	fmt.Println("Not enough money")
	return &pb.BuyingReply{Message: "Not enough money"}
}

func AccountIsBlocked(accountName string, accountId int) *pb.BuyingReply {
	fmt.Println("Operation is not permitted. Account is blocked. Name -", accountName, "ID -", accountId)
	return &pb.BuyingReply{Message: "Operation is not permitted. Account is blocked"}
}
