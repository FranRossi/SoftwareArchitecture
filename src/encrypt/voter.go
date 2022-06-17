package encrypt

import (
	"electoral_service/models"
)

func EncryptVoter(voter *models.VoterModel) {
	applyFunToVoter(voter, EncryptText)

}

func DecryptVoter(voter *models.VoterModel) {
	applyFunToVoter(voter, DecryptText)
}

func applyFunToVoter(voter *models.VoterModel, fun func(string) string) {
	voter.BirthDate = fun(voter.BirthDate)
	voter.Email = fun(voter.Email)
	voter.FullName = fun(voter.FullName)
	voter.LastCandidateVotedId = fun(voter.LastCandidateVotedId)
	voter.Phone = fun(voter.Phone)
	voter.Sex = fun(voter.Sex)

	for key, element := range voter.OtherFields {
		switch castedElement := element.(type) {
		case string:
			voter.OtherFields[key] = fun(castedElement)
		default:
			// Only string are encrypted or decrypted at the moment
			voter.OtherFields[key] = castedElement
		}

	}
}
