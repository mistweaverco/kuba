#!/usr/bin/env bash

# Needs to be run in macOS environment
# Imports Apple signing certificate and provisioning profile from GitHub secrets
# Creates temporary keychain and imports certificate into it
# Sets up outputs for use in later steps
#
# Requires the following secrets to be set and assigned to env in the repository:
# - BUILD_CERTIFICATE_BASE64: Base64 encoded .p12 certificate file
# - P12_PASSWORD: Base64 encoded password for the .p12 certificate
# - BUILD_PROVISION_PROFILE_BASE64: Base64 encoded .mobileprovision file
# - AUTH_KEY_BASE64: Base64 encoded .p8 auth key file for notarization

# create variables
CERTIFICATE_PATH=$RUNNER_TEMP/build_certificate.p12
PP_PATH=$RUNNER_TEMP/build_pp.provisionprofile
KEYCHAIN_PATH=$RUNNER_TEMP/app-signing.keychain-db
AUTH_KEY_PATH=$RUNNER_TEMP/AuthKey.p8
KEYCHAIN_PASSWORD=$(date +%s | sha256sum | base64 | head -c 32)
P12_PASSWORD=$(echo -n "$P12_PASSWORD" | base64 --decode)

# import certificate and provisioning profile from secrets
echo -n "$BUILD_CERTIFICATE_BASE64" | base64 --decode -o "$CERTIFICATE_PATH"
echo -n "$BUILD_PROVISION_PROFILE_BASE64" | base64 --decode -o "$PP_PATH"

# create temporary keychain
security create-keychain -p "$KEYCHAIN_PASSWORD" "$KEYCHAIN_PATH"
security set-keychain-settings -lut 21600 "$KEYCHAIN_PATH"
security unlock-keychain -p "$KEYCHAIN_PASSWORD" "$KEYCHAIN_PATH"

# import certificate to keychain
security import "$CERTIFICATE_PATH" -P "$P12_PASSWORD" -A -t cert -f pkcs12 -k "$KEYCHAIN_PATH"
security set-key-partition-list -S apple-tool:,apple: -k "$KEYCHAIN_PASSWORD" "$KEYCHAIN_PATH"
security list-keychain -d user -s "$KEYCHAIN_PATH"

# apply provisioning profile
mkdir -p ~/Library/MobileDevice/Provisioning\ Profiles
cp "$PP_PATH" ~/Library/MobileDevice/Provisioning\ Profiles

# create auth key file for notarization
echo -n "$AUTH_KEY_BASE64" | base64 --decode -o "$AUTH_KEY_PATH"

# setup outputs
{
  echo "auth_key_path=$AUTH_KEY_PATH";
  echo "keychain_path=$KEYCHAIN_PATH";
  echo "pp_path=$PP_PATH";
  echo "certificate_path=$CERTIFICATE_PATH";
} >> "$GITHUB_OUTPUT"
