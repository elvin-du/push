# register

Android首次Connect推送服务器Gate时，首先需要调用/register API生成一个RegID，AppID和RegID一起唯一标识一个用户。
IOS也需要调用/register API生成RegID，而且请求/register　API时，必须填写dev_token参数，Android则不需要，appID和dev_token唯一标识一个RegID.