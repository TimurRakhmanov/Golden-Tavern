# Golden Tavern
This is the repository for my bookings and reservation projects named Golden Tavern
<br/>A short summary of the project:
- Built in Go 1.20
- Uses the [chi router](https://github.com/go-chi/chi)
- Uses [alex edwards SCS](https://github.com/alexedwards/scs/v2) session management
- Uses [nosurf](https://github.com/justinas/nosurf)

# About the Project
The project itself is a part of my journey to learn Go and Backend development. 
There is a following functionality implemented:
* Search for days available for reservations for all rooms
* Choose a room out of available rooms and check its availability
* Inform the user when the reservation cannot be made
* If the reservation can be made, ask for further details
* If all necessary entries are verified and valid, show the summary of the reservation
* Send the email to both owner and the user about the reservation

# How Is It Implemented?
The project consists of the following elements:
* Backend written entirely in Go
* PostgreSQL is used as the database
* BootStrap framework is used for the front-end
* Third-party packages are used for form validation, email handling, session management, etc.
