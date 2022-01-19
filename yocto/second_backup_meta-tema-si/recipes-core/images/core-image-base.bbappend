
IMAGE_INSTALL += " dhcp-client go avahi-daemon ssh example python3"

inherit extrausers

EXTRA_USERS_PARAMS += "usermod -P labsi root;"
EXTRA_USERS_PARAMS += "useradd -P student student;"


