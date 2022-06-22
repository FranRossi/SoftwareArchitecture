package encrypt

type VoteModel struct {
	IdElection  string
	IdVoter     string
	Circuit     string
	IdCandidate string
	Signature   []byte
}

func EncryptVote(vote *VoteModel) {
	applyFunToVote(vote, EncryptText)
}

func DecryptVote(vote *VoteModel) {
	applyFunToVote(vote, DecryptText)
}

func applyFunToVote(vote *VoteModel, fun func(string) string) {
	vote.IdElection = fun(vote.IdElection)
	vote.IdVoter = fun(vote.IdVoter)
	vote.Circuit = fun(vote.Circuit)
	vote.IdCandidate = fun(vote.IdCandidate)
}
