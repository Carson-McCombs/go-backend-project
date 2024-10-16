Preplanning:

Development Choices:

Programming Language: Go / Golang

Reasoning: I have prior experience with multiple languages / web-frameworks that could be used here ( such as Go, PHP, Django / Python ), but I know Go is built for web services, scales well, and is optimized for concurrency. I also know that Fetch uses Go for some of its backend servers and services.

Development Environment: Visual Studio Code

Reasoning: It is open source, has plenty of support for many languages and frameworks, and I am used to it.

Other Tools: Docker

Reasoning: I've used Docker before and it works well for ensuring a quick deployment, testing, and compatibility. It also means that the user does not have to worry about setting anything else up or troubleshooting dependencies since it is run in pre-configured internal VM.

Libraries / Dependencies: Gorilla Mux

Reasoning: I have worked with Gorilla Mux before and know it is used within the industry as an HTTP router for Go web servers.

Abstraction Choices:

I could seperate the "Account" data into its own object, this would store the balance and transactions ( most as an array or list of Transactions within a Map, grouped by the Payer's name ) Pros:

*Consistent with the OOP approach

*Readability

*Scalability - even though the scope of the project was only set to a single account, by going ahead classifying this Account object, we can more easily scale the program to use more than one Account.

Cons:

*Not necessary for scope

Decision: Separate data into an "Account" object

Reasoning: For me, the next logical change for the program would be to go from one Account to multiple Accounts. And since I am planning the structure of my code in advance, I can go ahead and make sure that the code base will support such additions and we don't end up with extra refactoring or messy code down the line.

I could separate the "Transaction" data into its own object, this would store the Payer's name, the affect on points, and the Timestamp of the Transaction.

Pros:

*Consistent with the OOP approach

*The prompt explicitly used a Transaction object.

Cons:

*None really, any real issues with separating this data into its own object are more dealing with the nature of OOP than this instance.

Decision: Separate data into a "Transaction object"

Reasoning: Pretty straight forward choice will little negatives.

I could separate the "Balance" into its own object, which would store the Payer's name and the list of transactions for the respective account.

Pros:

*Consistent with the OOP approach

*Classifies related data as its own object ( perhaps cleaner? ) - this includes both store transactions and performing input validation with error handling.

*Perhaps could be better structured for future implementations - for example, if there was any other data being stored such as account tiers, we could store that data here.

Cons:

*Redundancy. A map would still be grouping Transactions by the Payer's name, it would just be in an array of Transactions.

*Data Redundancy / Duplication. This would also mean that there is a possible point of failure due to referencing the Payer's name in multiple areas. However, this is alrfeady the case to some degree due to a Transaction storing a reference to the Payer's name.

Decision: Do not separate data into a "Balance" object.

Reasoning: If we had to store any other data, such as balance per Payer or maybe if there was an Account tier associated with each Account per Payer, I would choose to separate the data into its own object. However, I will not separate the data into its own obejct, as that train of thought could steer towards over abstracting and classifying objects. As this can lead to a much larger amount of foundational work and complexity to get the program to run and maintain in the future.

I could group all of the relavent data into a "Server" object.

Pros:

*Consistent with the OOP approach

*Cleaner, no loose globally scoped objects.

*POST and GET API requests could be directly passed to the Server object to be handled

*Necessary data for loading and saving backups is already separated out

*More standard for creating Servers with Go

Cons:

*Similar to before, there aren't many cons to separating this data into its own object.

Decision: Group data into a Server object. Reasoning: Mainly for consistency, readability, and presedence, I've decided to separate / group the data into its own object.

Planned Schema:

Server extends mux.Router {

account Account

}

Account {

transactionsQueue map[string] ( Queue )

balance uint64

}

I used an unsigned 64-bit integer since the balance should never be negative and looks as though it is always a whole number.

I've decided to use the Payer's name as a key., considering the scope of the assignment, I would like to try and lower complexity and extra bloat as much as possible. But, it is primarily due to the POST requests always only reference the Payer's name, not a int or uid object. Therefore, I can only conclude that the Payer's name is considered an unique key.

Transaction {

payer string

points int64

timestamp Time

}

The Payer's name is a string.

The Transaction schema is basically directly from the project prompt. I decided to use a 64-bit integer to store the points because it can be either positive or negative and it appears to always be in whole numbers.

The timestamp uses Go's builtin "Time" standard library. I have done my own implementation for a timestamp on my "go-receipt-processor" project on my GitHub if you wish to see one, but as it is apart of the standard libraries, it has more first party support, testing, and updates than the alternative -as well as prevents me from introducing another point of errors or writing unecessary code.
