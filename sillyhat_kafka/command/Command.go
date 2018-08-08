package command

type CommandDTO struct {

	CommandName string `json:"commandName"`

	CommandBody string `json:"commandBody"`
}

func NewCommandDTO(commandName,commandBody string) *CommandDTO {
	return &CommandDTO{CommandName:commandName,CommandBody:commandBody}
}
