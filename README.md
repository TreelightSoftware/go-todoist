# Todist API

A simple unofficial Todoist REST API Client

Why pointers for the fields?

The default values have meaning in the Todoist API. In otherwords, if you try to update a task and set the content, but not the description field,
the description will default to "", so it will be set to blank without the user's intention. For further inspiration, see the excellent [Stripe-Go](https://github.com/stripe/stripe-go/) library.
