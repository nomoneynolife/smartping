package funcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// WechatRobotMessage 企业微信机器人消息结构
type WechatRobotMessage struct {
	MsgType  string                    `json:"msgtype"`
	Text     *WechatRobotTextContent   `json:"text,omitempty"`
	Markdown *WechatRobotMarkdownContent `json:"markdown,omitempty"`
}

// WechatRobotTextContent 文本消息内容
type WechatRobotTextContent struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

// WechatRobotMarkdownContent Markdown消息内容
type WechatRobotMarkdownContent struct {
	Content string `json:"content"`
}

// WechatRobotResponse 企业微信机器人响应
type WechatRobotResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// SendWechatRobotMessage 发送企业微信机器人消息
func SendWechatRobotMessage(webhookURL string, mentionedList string, mentionedMobile string, title string, content string) error {
	// 解析提醒列表
	var mentionedListArr []string
	if mentionedList != "" {
		mentionedListArr = strings.Split(mentionedList, ";")
	}
	
	var mentionedMobileArr []string
	if mentionedMobile != "" {
		mentionedMobileArr = strings.Split(mentionedMobile, ";")
	}
	
	// 构建消息内容
	message := WechatRobotMessage{
		MsgType: "text",
		Text: &WechatRobotTextContent{
			Content:             fmt.Sprintf("%s\n\n%s", title, content),
			MentionedList:       mentionedListArr,
			MentionedMobileList: mentionedMobileArr,
		},
	}
	
	// 转换为JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}
	
	// 发送HTTP请求
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		return fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}
	
	// 解析响应
	var response WechatRobotResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("JSON解析失败: %v", err)
	}
	
	// 检查响应状态
	if response.ErrCode != 0 {
		return fmt.Errorf("企业微信机器人返回错误: %s (错误码: %d)", response.ErrMsg, response.ErrCode)
	}
	
	return nil
}

// SendWechatRobotAlert 发送企业微信机器人告警消息
func SendWechatRobotAlert(alertConfig map[string]string, alertTitle string, alertContent string) error {
	// 获取配置
	webhookURL := alertConfig["WechatWebhook"]
	mentionedList := alertConfig["WechatMentionedList"]
	mentionedMobile := alertConfig["WechatMentionedMobile"]
	
	// 检查Webhook地址是否配置
	if webhookURL == "" {
		return fmt.Errorf("企业微信机器人Webhook地址未配置")
	}
	
	// 发送消息
	err := SendWechatRobotMessage(webhookURL, mentionedList, mentionedMobile, alertTitle, alertContent)
	if err != nil {
		return fmt.Errorf("发送企业微信机器人消息失败: %v", err)
	}
	
	return nil
}

// SendWechatRobotTest 发送企业微信机器人测试消息
func SendWechatRobotTest(webhookURL string, mentionedList string, mentionedMobile string) error {
	testTitle := "SmartPing 测试消息"
	testContent := "这是一条来自 SmartPing 的测试消息，用于验证企业微信机器人配置是否正确。"
	
	err := SendWechatRobotMessage(webhookURL, mentionedList, mentionedMobile, testTitle, testContent)
	if err != nil {
		return fmt.Errorf("发送测试消息失败: %v", err)
	}
	
	return nil
}