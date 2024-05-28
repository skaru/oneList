# oneList

This is a simple to-do list in Go. The idea was to keep it as simple as possible and make it easy to deploy. The binary contains its own webserver, website, storage and authentication. This was my first attempt at OOP Go. I primarily used interfaces to make it easy to extend with different frontend and storage methods.

There were three requirements:
- Must easily sync between devices
- Must be easy to install
- Must not need an account/login

Making the list accessible as website already fulfils the first two requirements. However, making it a website means that everyone could see my personal to-do list. So I added a check for a certain cookie before any content is loaded. The cookie can be set by visiting a special page on which a password has to be entered. Although the third requirement isn't met by doing this. I consider it to be sufficient because the login is only required once per device.
