package common

import "log"

func HandleError(err error) error {
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
