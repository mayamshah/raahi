from flask import Flask, redirect, request
from stravalib.client import Client
import requests
import os
import time

CLIENTID = 50880
CLIENTSECRET = "179f538ea081f2553c200441892e8fde3dc5255e"

api = Flask(__name__)

CLIENTID = 50880
CLIENTSECRET = "179f538ea081f2553c200441892e8fde3dc5255e"
client = Client()


@api.route("/")
def root():
    return("Welcome to strava-oauth")


@api.route("/authorize")
def authorize():
    """Redirect user to the Strava Authorization page"""
    authorize_url = client.authorization_url(client_id=CLIENTID, redirect_uri='http://localhost:5000/authorization_successful')
    return redirect(authorize_url)


@api.route("/authorization_successful")
def authorization_successful():
    """Exchange code for a user token"""
    params = {
        "client_id": CLIENTID,
        "client_secret": CLIENTSECRET,
        "code": request.args.get('code'),
        "grant_type": "authorization_code"
    }

    r = requests.post("https://www.strava.com/oauth/token", params)

    client.token_expires_at = r.json()["expires_at"]
    client.access_token = r.json()["access_token"]
    client.refresh_token = r.json()["refresh_token"]
    return "Authorization successful"


@api.route("/token") 
def token():
    if time.time() > client.token_expires_at:
        refresh_response = client.refresh_access_token(client_id=CLIENTID, client_secret=CLIENTSECRET,
            refresh_token=client.refresh_token)
        access_token = refresh_response['access_token']
        refresh_token = refresh_response['refresh_token']
        expires_at = refresh_response['expires_at']
        client.token_expires_at = token_expires_at
        client.access_token = access_token
        client.refresh_token - refresh_token
    else:
        print(client.access_token)

    return "{\"accessToken\": \"" + client.access_token + "\"}"

if __name__ == '__main__':
    api.run()