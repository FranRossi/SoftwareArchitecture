package main

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"electoral_service/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	math "math/rand"
	l "own_logger"
	"strconv"
	"sync"
	"voting_simulator/proto"
)

func GetVoters() ([]models.VoterModel, error) {
	client := GetInstanceMongoClient()
	collection := client.Database("uruguay_election").Collection("voters")
	var results []bson.M
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		l.LogError(err.Error())
	}
	defer cursor.Close(context.TODO())
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
	if err != nil {
		l.LogError(err.Error())
	}
	fmt.Println("Started voting")
	wg := sync.WaitGroup{}
	wg.Add(len(voters))
	for _, voter := range voters {
		go CreateVote(voter, &wg)
	}
	wg.Wait()
	fmt.Println("Finished voting")
}

func CreateVote(voter models.VoterModel, wg *sync.WaitGroup) {
	defer wg.Done()
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

const addr = ":50004"

func Vote(vote VoteModel) {
	encryptVote(&vote)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot connect: %s", err)
	}
	defer conn.Close()

	client := proto.NewVoteServiceClient(conn)
	request := &proto.VoteRequest{
		IdElection:  vote.IdElection,
		IdVoter:     vote.IdVoter,
		Circuit:     vote.Circuit,
		IdCandidate: vote.IdCandidate,
		Signature:   vote.Signature,
	}
	response, err2 := client.Vote(context.Background(), request)
	if err2 != nil {
		log.Fatalf("could not vote: %v", err2)
	}
	fmt.Println("Vote: %s\n", response.Message)
}
