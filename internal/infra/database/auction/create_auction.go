package auction

import (
	"context"
	"fmt"
	"os"
	"time"

	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(ctx context.Context, auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error inserting auction", err)
		return internal_error.NewInternalServerError("Error inserting auction")
	}

	return nil
}

func (ar *AuctionRepository) StartAuctionCloser(ctx context.Context) {
	checkInterval := getCheckInterval()
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	logger.Info(fmt.Sprintf("Auction closer started, checking every %s", checkInterval))

	for {
		select {
		case <-ticker.C:
			ar.closeExpiredAuctions(ctx)
		case <-ctx.Done():
			logger.Info("Auction closer stopped.")
			return
		}
	}
}

func (ar *AuctionRepository) closeExpiredAuctions(ctx context.Context) {
	now := time.Now().Unix()
	auctionDuration := getAuctionInterval()

	filter := bson.M{
		"status":    auction_entity.Active,
		"timestamp": bson.M{"$lt": now - int64(auctionDuration.Seconds())},
	}
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

	result, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error closing expired auctions:", err)
		return
	}

	if result.ModifiedCount > 0 {
		logger.Info(fmt.Sprintf("Closed %d expired auctions.", result.ModifiedCount))
	}
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	if auctionInterval == "" {
		return 5 * time.Minute
	}
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return 5 * time.Minute
	}
	return duration
}

func getCheckInterval() time.Duration {
	checkInterval := os.Getenv("AUCTION_CHECK_INTERVAL")
	if checkInterval == "" {
		return 30 * time.Second
	}
	duration, err := time.ParseDuration(checkInterval)
	if err != nil {
		return 30 * time.Second
	}
	return duration
}
