package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/entity/auction_entity"

	"go.mongodb.org/mongo-driver/bson"
)

func TestCloseExpiredAuctions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	os.Setenv("MONGODB_URL", "mongodb://admin:admin@localhost:27017/auctions?authSource=admin")
	os.Setenv("MONGODB_DB", "auctions")
	os.Setenv("AUCTION_INTERVAL", "20s")
	defer func() {
		os.Unsetenv("MONGODB_URL")
		os.Unsetenv("MONGODB_DB")
		os.Unsetenv("AUCTION_INTERVAL")
	}()

	dbConn, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v.", err)
	}
	defer dbConn.Client().Disconnect(ctx)

	testDBName := dbConn.Name() + "_" + time.Now().Format("20060102150405")
	testDB := dbConn.Client().Database(testDBName)
	defer func() {
		if err := testDB.Drop(ctx); err != nil {
			t.Errorf("Failed to drop test database: %v", err)
		}
	}()

	repo := NewAuctionRepository(testDB)

	auction := &auction_entity.Auction{
		Id:          "test-auction",
		ProductName: "Produto Teste",
		Category:    "Categoria Teste",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now().Add(-25 * time.Second),
	}

	if err := repo.CreateAuction(ctx, auction); err != nil {
		t.Fatalf("Failed to create auction: %v", err)
	}

	repo.closeExpiredAuctions(ctx)

	var closedAuction AuctionEntityMongo
	err = repo.Collection.FindOne(ctx, bson.M{"_id": auction.Id}).Decode(&closedAuction)
	if err != nil {
		t.Fatalf("Failed to find auction: %v", err)
	}

	if closedAuction.Status != auction_entity.Completed {
		t.Errorf("Expected auction status to be %v, got %v", auction_entity.Completed, closedAuction.Status)
	}
}
