/* Copyright 2018 sky<skygangsta@hotmail.com>. All rights reserved.
 *
 * Licensed under the Apache License, version 2.0 (the "License").
 * You may not use this work except in compliance with the License, which is
 * available at www.apache.org/licenses/LICENSE-2.0
 *
 * This software is distributed on an "AS IS" basis, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied, as more fully set forth in the License.
 *
 * See the NOTICE file distributed with this work for information regarding copyright ownership.
 */

package logger

type Color struct{}

func NewColor() *Color {
	return &Color{}
}

func (this *Color) White() string {
	return "\033[30m"
}

func (this *Color) Red() string {
	return "\033[31m"
}

func (this *Color) Green() string {
	return "\033[32m"
}

func (this *Color) Yello() string {
	return "\033[33m"
}

func (this *Color) Blue() string {
	return "\033[34m"
}

func (this *Color) Magenta() string {
	return "\033[35m"
}

func (this *Color) Cyan() string {
	return "\033[36m"
}

func (this *Color) Clear() string {
	return "\033[0m"
}

func (this *Color) WhiteBackground() string {
	return "\033[97;40m"
}

func (this *Color) RedBackground() string {
	return "\033[97;41m"
}

func (this *Color) GreenBackground() string {
	return "\033[97;42m"
}

func (this *Color) YelloBackground() string {
	return "\033[97;43m"
}

func (this *Color) BlueBackground() string {
	return "\033[97;44m"
}

func (this *Color) MagentaBackground() string {
	return "\033[97;45m"
}

func (this *Color) CyanBackground() string {
	return "\033[97;46m"
}

// \033[90;47m
func (this *Color) GrayBackground() string {
	return "\033[97;47m"
}
