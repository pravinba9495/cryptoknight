FROM node:16-alpine3.14

# Installs latest Chromium package.
RUN apk add --no-cache \
      chromium \
      nss \
      freetype \
      harfbuzz \
      ca-certificates \
      ttf-freefont

# Tell Puppeteer to skip installing Chrome. We'll be using the installed package.
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true \
    PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

# Add user so we don't need --no-sandbox.
RUN addgroup -S pptruser && adduser -S -g pptruser pptruser \
    && mkdir -p /home/pptruser/Downloads /app \
    && chown -R pptruser:pptruser /home/pptruser \
    && chown -R pptruser:pptruser /app

# Run everything after as non-privileged user.
USER pptruser

RUN mkdir -p /home/pptruser/pravinba9495/kryptonite
COPY dist /home/pptruser/pravinba9495/kryptonite
COPY package.json /home/pptruser/pravinba9495/kryptonite
COPY package-lock.json /home/pptruser/pravinba9495/kryptonite
WORKDIR /home/pptruser/pravinba9495/kryptonite

USER root
RUN npm install

USER pptruser

CMD node index.js