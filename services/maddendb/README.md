# Madden Database library 

This library handles all configuration and CRUD actions for the Madden database. 

## Configuration 

Configuration items as detailed in [standard database configuration](../dbutils/README.md)

## Testing 

This library expects a PostgreSQL database to be available for testing. The supplied script runtests.sh can be used to set env and launch the tests once you have a database available. This script will set appropriate environment variables, but note they may need to be updated for your environment. If attempting to run tests manually, note these variables and set them accordingly. 