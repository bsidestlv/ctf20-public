import os
import boto3


def lambda_handler(event, context):
    client = client = boto3.client('cognito-idp')
    try:
        client.admin_get_user(
            UserPoolId=os.environ['UserPoolId'],
            Username='marx'
        )
    except:
        client.admin_create_user(
            UserPoolId=os.environ['UserPoolId'],
            Username='marx',
            UserAttributes=[
                {
                    'Name': 'custom:isWarden',
                    'Value': '0'
                },
            ]
        )
    client.admin_set_user_password(
        UserPoolId=os.environ['UserPoolId'],
        Username='marx',
        Password='Jailbreak',
        Permanent=True
    )
    return event