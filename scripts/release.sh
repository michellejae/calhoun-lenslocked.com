#!/bin/bash
cd "$GOPATH/src/lenslocked.com"
echo "==== Releasing lenslocked.com ===="
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm lenslocked.com
echo "  Done!"
echo "  Deleting existing code..."
ssh root@want2breakfree.com "rm -rf /root/go/src/lenslocked.com"
echo "  Code deleted successfully!"
echo "  Uploading code..."
# The \ at the end of the line tells bash that our
# command isn't done and wraps to the next line.
rsync -avr --exclude '.git/*' --exclude 'tmp/*' \
  --exclude 'images/*' ./ \
  root@want2breakfree.com:/root/go/src/lenslocked.com/
echo "  Code uploaded successfully!"
echo "  Go getting deps..."
ssh root@want2breakfree.com "export GOPATH=/root/go; \
  /usr/local/go/bin/go get golang.org/x/crypto/bcrypt"
ssh root@want2breakfree.com "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/mux"
ssh root@want2breakfree.com "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/schema"
ssh root@want2breakfree.com "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/lib/pq"
ssh root@want2breakfree.com "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/jinzhu/gorm"
ssh root@want2breakfree.com "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/csrf"
ssh root@want2breakfree.com "export GOPATH=/root/go; \
  /usr/local/go/bin/go get gopkg.in/mailgun/mailgun-go.v1"

echo "  Building the code on remote server..."
ssh root@want2breakfree.com 'cd /root/go/src/lenslocked.com; export GOPATH=/root/go; /usr/local/go/bin/go build -o /root/app/server $GOPATH/src/lenslocked.com/*.go'
echo "  Code built successfully!"

echo "  Moving assets..."
ssh root@want2breakfree.com "cd /root/app; \
  cp -R /root/go/src/lenslocked.com/assets ."
echo "  Assets moved successfully!"

echo "  Moving views..."
ssh root@want2breakfree.com "cd /root/app; \
  cp -R /root/go/src/lenslocked.com/views ."
echo "  Views moved successfully!"

echo "  Moving Caddyfile..."
ssh root@want2breakfree.com "cp /root/go/src/lenslocked.com/Caddyfile /etc/caddy/Caddyfile"
echo "  Caddyfile moved successfully!"

echo "  Restarting the server..."
ssh root@want2breakfree.com "sudo service lenslocked.com restart"
echo "  Server restarted successfully!"

echo "  Restarting Caddy server..."
ssh root@want2breakfree.com "sudo service caddy restart"
echo "  Caddy restarted successfully!"

echo "==== Done releasing lenslocked.com ===="