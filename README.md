# Yahoo
Go package for consuming Yahoo! fantasy sports API.

The top level package contains one small Go file auth.go
This provides a pre-configured golang.org/x/oauth2.Config for Yahoo! connections.
Dependent on golang.org/x/oauth2 package

 **NewConfig** creates a ready to use oauth2.Config instance. if the redirectUrl is empty,
 this is assumed to be an installed application.
 
`func NewConfig(clientId, clientSecret, redirectUrl string) *oauth2.Config`

**Example**

`conf := yahoo.NewConfig("MyYahooClientID", "MyYahooClientSecret", "http://mysite.com/callback")`

## Fantasy
The fantasy sub-package contains structs, methods and query builders to consume the fantasy API. 
Fantasy has no dependencies outside the standard library.

**WARNING - This code is in its infancy and the api is likely to change frequently.** 
