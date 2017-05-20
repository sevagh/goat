package main

func Mkdir(mountPath string) error {
	cmd := "mkdir"
	args := []string{
		"-p",
		mountPath,
	}
	if _, err := ExecuteCommand(cmd, args); err != nil {
		return err
	}
	return nil
}
