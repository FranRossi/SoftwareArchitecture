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
		// case int:
		// voter.OtherFields[key], _ = strconv.ParseInt(fun(strconv.Itoa(castedElement)), 10, 64)
		// voter.OtherFields[key] = fun(string((castedElement)))
		//fun(string([]byte(castedElement)))
		// case bool:
		// voter.OtherFields[key], _ = strconv.ParseBool(fun(strconv.FormatBool(castedElement)))
		// case float64:
		// voter.OtherFields[key], _ = strconv.ParseFloat(fun(strconv.FormatFloat(castedElement, 'E', -1, 64)), 64)
		// case uint64:
		// voter.OtherFields[key], _ = strconv.ParseUint(fun(strconv.FormatUint(castedElement, 10)), 10, 64)
		default:
			// Some less common types are not handled yet
			// and are not encrypted or decrypted (e.g. arrays or maps that could cause performance issues)
			voter.OtherFields[key] = castedElement
		}

	}
}
