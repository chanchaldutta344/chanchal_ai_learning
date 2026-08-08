Feature: Customer API
  As a bank application user
  I want to create, query, and update customer records
  So that customer profile data can be stored and managed through the API

  Scenario: Create a customer through the REST API
    Given the customer service is running
    When I POST a customer payload to /customers
    Then the response status should be 201
    And the customer should contain a customer_id

  Scenario: Read a customer by id through the REST API
    Given a customer exists in the customers table
    When I GET /customers/{customer_id}
    Then the response status should be 200
    And the JSON body should return the customer record

  Scenario: Update a customer through the REST API
    Given a customer exists in the customers table
    When I PUT a modified customer payload to /customers/{customer_id}
    Then the response status should be 200
    And the JSON body should contain the updated email field
