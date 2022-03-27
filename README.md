Version 1.1

Brief:
Design an API for a vending machine, allowing users with a “seller” role to add, update or remove products, while users
with a “buyer” role can deposit coins into the machine and make purchases. Your vending machine should only accept 5,
10, 20, 50 and 100 cent coins.

Tasks:
REST API should be implemented consuming and producing “application/json” Implement product model with amountAvailable,
cost, productName and sellerId fields Implement user model with username, password, deposit and role fields Implement

CRUD for users (POST shouldn’t require authentication)

Implement CRUD for a product model (GET can be called by anyone, while POST, PUT and DELETE can be called only by the
seller user who created the product)

Implement /deposit endpoint so users with a “buyer” role can deposit 5, 10, 20, 50 and 100 cent coins into their vending
machine account Implement /buy endpoint (accepts productId, amount of products) so users with a “buyer” role can buy

products with the money they’ve deposited. API should return total they’ve spent, products they’ve purchased and their
change if there’s any (in 5, 10, 20, 50 and 100 cent coins),

v1.1 Software now take supported coins as input. Default is still 5,10,20,50,100

Implement /reset endpoint so users with a “buyer” role can reset their deposit

### Run project:

docker-compose up -d

### Stop project:

docker-compose down

### Test:

```make test```

## Integration testing(POSTGRESQL):

```make test-integration```

## TODO:

Add API documentation

Replace chi router with echo framework since it offer a nicer error handling where the errors are returned and could be
handled by a middleware, Or consider passing the errors down by context and resolving in a middleware


Instead of seperating sql code in a private method and wrap it with error handling, 
Could let codegenerating generate the functions from sql files and call the generated function directly in the repository wrapper 


## DONE:

Write a ci/cd pipeline   
