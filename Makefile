build:
	go build -o all-the-heics.app/Contents/MacOS/all-the-heics .

build-icons:
	iconutil -c icns -o icon.icns ath.iconset
