package main

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"electoral_service/models"
	"fmt"
	math "math/rand"
	"os"
	l "own_logger"
	"strconv"
	"sync"
	"time"
	"voting_simulator/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetVoters() ([]models.VoterModel, error) {
	client := GetInstanceMongoClient()
	collection := client.Database("uruguay_election").Collection("voters")
	var results []bson.M
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		l.LogError(err.Error())
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			l.LogError(err.Error())
		}
	}(cursor, context.TODO())
	if err = cursor.All(context.TODO(), &results); err != nil {
		return []models.VoterModel{}, fmt.Errorf("error requesting voters: %v", err)
	}
	voters := convertResultsToVoters(results)
	return voters, nil
}

func convertResultsToVoters(results []bson.M) []models.VoterModel {
	voters := []models.VoterModel{}
	for _, result := range results {
		other := result["otherFields"].(bson.M)
		user := &models.VoterModel{
			Id:                   result["id"].(string),
			FullName:             result["name"].(string),
			Sex:                  result["sex"].(string),
			BirthDate:            result["birthDate"].(string),
			Phone:                result["phone"].(string),
			Email:                result["email"].(string),
			Voted:                int(result["voted"].(int32)),
			LastCandidateVotedId: result["lastCandidateVotedId"].(string),
			Region:               result["region"].(string),
			OtherFields:          other,
		}
		voters = append(voters, *user)
	}
	fmt.Println("Total voters: " + strconv.Itoa(len(voters)))
	return voters
}

func CreateVotes() {
	voters, err := GetVoters()
	limit := os.Getenv("LIMIT_VOTES")
	limitVotes, err := strconv.Atoi(limit)
	if err != nil {
		l.LogError(err.Error())
	}
	fmt.Println("Started voting")
	startTime := time.Now()
	var wg sync.WaitGroup
	count := 0
	batchString := os.Getenv("BATCH_SIZE")
	batchSize, _ := strconv.Atoi(batchString)
	fmt.Println("Batch-size: " + batchString)
	for i := 0; i < len(voters) && i <= limitVotes; i++ {
		localVoter := voters[i]
		count++
		wg.Add(1)
		go func(v models.VoterModel) {
			CreateVote(localVoter)
			wg.Done()
		}(localVoter)
		if count >= batchSize {
			wg.Wait()
			count = 0
		}
	}
	wg.Wait()
	endTime := time.Now()
	fmt.Println("Finished voting")
	fmt.Println("Time: ", endTime.Sub(startTime))
}

func CreateVote(voter models.VoterModel) {
	candidateId := strconv.Itoa(int(math.Int31n(3) + 1))
	vote := VoteModel{
		IdElection:  "1",
		IdVoter:     voter.Id,
		Circuit:     voter.OtherFields["circuit"].(string),
		IdCandidate: candidateId,
	}
	SignAndVote(vote)
}

func SignAndVote(vote VoteModel) {
	privateKey := getPrivateKey()
	voter := []byte(vote.IdVoter)
	msgHash := sha256.New()
	msgHash.Write(voter)
	msgHashSBytes := msgHash.Sum(nil)
	signature, _ := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSBytes, nil)
	vote.Signature = signature
	Vote(vote)
}

type VoteModel struct {
	IdElection  string
	IdVoter     string
	Circuit     string
	IdCandidate string
	Signature   []byte
}

func Vote(vote VoteModel) {
	addr := os.Getenv("GRPC")
	encryptVote(&vote)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("cannot connect: %s", err)
		conn.Close()
		return
	}

	client := proto.NewVoteServiceClient(conn)
	request := &proto.VoteRequest{
		IdElection:  vote.IdElection,
		IdVoter:     vote.IdVoter,
		Circuit:     vote.Circuit,
		IdCandidate: vote.IdCandidate,
		Signature:   vote.Signature,
	}
	_, err2 := client.Vote(context.Background(), request)
	if err2 != nil {
		// fmt.Printf("\n could not vote: %v", err2)
		fmt.Printf("X")
		conn.Close()
		return
	}
	// fmt.Println("Vote: %s\n", response.Message)
	fmt.Print(".")
	conn.Close()
}
