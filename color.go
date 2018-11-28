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

//# 红色  echo -e "\033[31m[ERROR] \033[0m"
//# 绿色  echo -e "\033[32m \033[0m"
//# 黄色  echo -e "\033[33m[ALERT] \033[0m"
//# 蓝色  echo -e "\033[34m[UNUSED] \033[0m"
//# 粉色  echo -e "\033[35m[NOTICE] \033[0m"
//# 青色  echo -e "\033[36m[INFO] \033[0m"
//green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
//white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
//yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
//red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
//blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
//magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
//cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
//reset   = string([]byte{27, 91, 48, 109})

package logger

type Color struct {
}

func NewColor() *Color {
	return &Color{}
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

func (this *Color) RedBackground() string {
	return string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
}

func (this *Color) GreenBackground() string {
	return string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
}

func (this *Color) YelloBackground() string {
	return string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
}

func (this *Color) BlueBackground() string {
	return string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
}

func (this *Color) MagentaBackground() string {
	return string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
}

func (this *Color) CyanBackground() string {
	return string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
}

func (this *Color) WhiteBackground() string {
	return string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
}
