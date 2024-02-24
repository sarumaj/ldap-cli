Feature: Modify requests
    
    Scenario: I submit a modify request to change user's password
        Given I perform bind with following parameters:
            | --user               | cn=admin,dc=mock,dc=ad,dc=com |
            | --password           | admin                         |
            | --url                | ldap://localhost:389          |
        And I plan to execute "edit" command with following parameters:
            | --path               | dc=mock,dc=ad,dc=com          |
        And I plan to execute "user" command with following parameters:
            | --user-id            | uix00001                      |
            | --new-password       | new-password                  |
            | --password-attribute | userPassword                  |
        When I execute the application
        Then I expect the output to be:
            """
            Successfully applied modifications
            """
        And I expect the exit code to be 0

    Scenario: I submit a modify request to replace group's members
        Given I perform bind with following parameters:
            | --user             | cn=admin,dc=mock,dc=ad,dc=com |
            | --password         | admin                         |
            | --url              | ldap://localhost:389          |
        And I plan to execute "edit" command with following parameters:
            | --path             | dc=mock,dc=ad,dc=com          |
        And I plan to execute "group" command with following parameters:
            | --group-id         | group01                       |
            | --replace-member   | uix00002                      |
            | --member-attribute | memberUid                     |
        When I execute the application
        Then I expect the output to be:
            """
            Successfully applied modifications
            """
        And I expect the exit code to be 0