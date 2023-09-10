package objects

/*Contains information about Telegram Passport data shared with the bot by the user.*/
type PassportData struct {
	/*Array with information about documents and other Telegram Passport elements that was shared with the bot*/
	Data []EncryptedPassportElement `json:"data"`
	/*Encrypted credentials required to decrypt the data*/
	Credentials *EncryptedCredentials `json:"credentials"`
}

/*This object represents a file uploaded to Telegram Passport. Currently all Telegram Passport files are in JPEG format when decrypted and don't exceed 10MB.*/
type PassportFile struct {
	/*Identifier for this file, which can be used to download or reuse the file*/
	FileId string `json:"file_id"`
	/*Unique identifier for this file, which is supposed to be the same over time and for different bots. Can't be used to download or reuse the file.*/
	FileUniqueId string `json:"file_unique_id"`
	/*File size in bytes*/
	FileSize int `json:"file_size"`
	/*Unix time when the file was uploaded*/
	FileDate int `json:"file_date"`
}

/*Contains information about documents or other Telegram Passport elements shared with the bot by the user.*/
type EncryptedPassportElement struct {
	/*Element type. One of “personal_details”, “passport”, “driver_license”, “identity_card”, “internal_passport”, “address”, “utility_bill”, “bank_statement”, “rental_agreement”, “passport_registration”, “temporary_registration”, “phone_number”, “email”.*/
	Type string `json:"type"`
	/*Optional. Base64-encoded encrypted Telegram Passport element data provided by the user, available for “personal_details”, “passport”, “driver_license”, “identity_card”, “internal_passport” and “address” types. Can be decrypted and verified using the accompanying EncryptedCredentials.*/
	Data string `json:"data,omitempty"`
	/*Optional. User's verified phone number, available only for “phone_number” type*/
	PhoneNumber string `json:"phone_number,omitempty"`
	/*Optional. User's verified email address, available only for “email” type*/
	Email string `json:"email,omitempty"`
	/*Optional. Array of encrypted files with documents provided by the user, available for “utility_bill”, “bank_statement”, “rental_agreement”, “passport_registration” and “temporary_registration” types. Files can be decrypted and verified using the accompanying EncryptedCredentials*/
	Files []PassportFile `json:"files,omitempty"`
	/*Optional. Encrypted file with the front side of the document, provided by the user. Available for “passport”, “driver_license”, “identity_card” and “internal_passport”. The file can be decrypted and verified using the accompanying EncryptedCredentials.*/
	FrontSide *PassportFile `json:"front_side,omitempty"`
	/*Optional. Encrypted file with the reverse side of the document, provided by the user. Available for “driver_license” and “identity_card”. The file can be decrypted and verified using the accompanying EncryptedCredentials.*/
	ReverseSide *PassportFile `json:"reverse_side,omitempty"`
	/*Optional. Encrypted file with the selfie of the user holding a document, provided by the user; available for “passport”, “driver_license”, “identity_card” and “internal_passport”. The file can be decrypted and verified using the accompanying EncryptedCredentials.*/
	Selfie *PassportFile `json:"selfie,omitempty"`
	/*Optional. Array of encrypted files with translated versions of documents provided by the user. Available if requested for “passport”, “driver_license”, “identity_card”, “internal_passport”, “utility_bill”, “bank_statement”, “rental_agreement”, “passport_registration” and “temporary_registration” types. Files can be decrypted and verified using the accompanying EncryptedCredentials.*/
	Translation []PassportFile `json:"translation,omitempty"`
	/*	Base64-encoded element hash for using in PassportElementErrorUnspecified*/
	Hash string `json:"hash"`
}

/*Contains data required for decrypting and authenticating EncryptedPassportElement. See the Telegram Passport Documentation for a complete description of the data decryption and authentication processes.*/
type EncryptedCredentials struct {
	/*Base64-encoded encrypted JSON-serialized data with unique user's payload, data hashes and secrets required for EncryptedPassportElement decryption and authentication*/
	Data string `json:"data"`
	/*Base64-encoded data hash for data authentication*/
	Hash string `json:"hash"`
	/*Base64-encoded secret, encrypted with the bot's public RSA key, required for data decryption*/
	Secret string `json:"secret"`
}

type PassportElementError interface {
	blah()
}

/*This object represents an error in the Telegram Passport element which was submitted that should be resolved by the user. It should be one of:

PassportElementErrorDataField
PassportElementErrorFrontSide
PassportElementErrorReverseSide
PassportElementErrorSelfie
PassportElementErrorFile
PassportElementErrorFiles
PassportElementErrorTranslationFile
PassportElementErrorTranslationFiles
PassportElementErrorUnspecified

This object should not be used at all.*/
type PassportElementErrordefault struct {
	/*Error source, must be data*/
	Source string `json:"source"`
	/*The section of the user's Telegram Passport which has the error, one of “personal_details”, “passport”, “driver_license”, “identity_card”, “internal_passport”, “address”*/
	Type string `json:"type"`
	/*Base64-encoded hash of the file with the front side of the document*/
	FileHash string `json:"file_hash"`
	/*Error message*/
	Message string `json:"message"`
}

/*Represents an issue in one of the data fields that was provided by the user. The error is considered resolved when the field's value changes.*/
type PassportElementErrorDataField struct {
	/*Error source, must be data*/
	Source string `json:"source"`
	/*The section of the user's Telegram Passport which has the error, one of “personal_details”, “passport”, “driver_license”, “identity_card”, “internal_passport”, “address”*/
	Type string `json:"type"`
	/*Error message*/
	Message string `json:"message"`
	/*Name of the data field which has the error*/
	FieldName string `json:"field_name"`
	/*Base64-encoded data hash*/
	DataHash string `json:"data_hash"`
}

func (ps *PassportElementErrorDataField) blah() {}

/*Represents an issue with the front side of a document. The error is considered resolved when the file with the front side of the document changes.*/
type PassportElementErrorFrontSide struct {
	PassportElementErrordefault
}

func (ps *PassportElementErrorFrontSide) blah() {}

/*Represents an issue with the reverse side of a document. The error is considered resolved when the file with reverse side of the document changes.*/
type PassportElementErrorReverseSide struct {
	PassportElementErrordefault
}

func (ps *PassportElementErrorReverseSide) blah() {}

/*Represents an issue with the selfie with a document. The error is considered resolved when the file with the selfie changes.*/
type PassportElementErrorSelfie struct {
	PassportElementErrordefault
}

func (ps *PassportElementErrorSelfie) blah() {}

/*Represents an issue with a document scan. The error is considered resolved when the file with the document scan changes.*/
type PassportElementErrorFile struct {
	PassportElementErrordefault
}

func (ps *PassportElementErrorFile) blah() {}

/*Represents an issue with a list of scans. The error is considered resolved when the list of files containing the scans changes.*/
type PassportElementErrorFiles struct {
	/*Error source, must be data*/
	Source string `json:"source"`
	/*The section of the user's Telegram Passport which has the error, one of “personal_details”, “passport”, “driver_license”, “identity_card”, “internal_passport”, “address”*/
	Type string `json:"type"`
	/*List of base64-encoded file hashes*/
	FileHashes []string `json:"file_hashes"`
	/*Error message*/
	Message string `json:"message"`
}

func (ps *PassportElementErrorFiles) blah() {}

/*Represents an issue with one of the files that constitute the translation of a document. The error is considered resolved when the file changes.*/
type PassportElementErrorTranslationFile struct {
	PassportElementErrordefault
}

func (ps *PassportElementErrorTranslationFile) blah() {}

/*Represents an issue with the translated version of a document. The error is considered resolved when a file with the document translation change.*/
type PassportElementErrorTranslationFiles struct {
	PassportElementErrorFiles
}

func (ps *PassportElementErrorTranslationFiles) blah() {}

/*Represents an issue in an unspecified place. The error is considered resolved when new data is added.*/
type PassportElementErrorUnspecified struct {
	/*Error source, must be unspecified*/
	Source string `json:"source"`
	/*Type of element of the user's Telegram Passport which has the issue*/
	Type string `json:"type"`
	/*Base64-encoded element hash*/
	ElementHash string `json:"element_hash"`
	/*Error message*/
	Message string `json:"mesaage"`
}

func (ps *PassportElementErrorUnspecified) blah() {}
