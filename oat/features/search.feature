Feature: Search requests

    Scenario: I submit a search request to find a specific user using custom filter
        Given I perform bind with following parameters:
        | --user     | cn=admin,dc=mock,dc=ad,dc=com |
        | --password | admin                         |
        | --url      | ldap://localhost:389          |
        And I plan to execute "get" command with following parameters:
        | --path   | dc=mock,dc=ad,dc=com |
        | --select | *                    |
        And I plan to execute "custom" command with following parameters:
        | --filter | (cn=uix00001) |
        When I execute the application
        Then I expect the output to be:
        """
        CommonName: uix00001
        Gecos: uix00001
        GidNumber: \"100\"
        HomeDirectory: /home/uix00001
        LoginShell: /bin/bash
        ObjectClass:
        - top
        - account
        - posixAccount
        - shadowAccount
        ShadowLastChange: \"0\"
        ShadowMax: \"0\"
        ShadowWarning: \"0\"
        Uid: uix00001
        UidNumber: \"\d+\"
        UserPassword: new-password
        """
        And I expect the exit code to be 0

    Scenario: I submit a search request to find a specific group using custom filter
        Given I perform bind with following parameters:
        | --user     | cn=admin,dc=mock,dc=ad,dc=com |
        | --password | admin                         |
        | --url      | ldap://localhost:389          |
        And I plan to execute "get" command with following parameters:
        | --path   | dc=mock,dc=ad,dc=com |
        | --select | *                    |
        And I plan to execute "custom" command with following parameters:
        | --filter | (cn=group01) |
        When I execute the application
        Then I expect the output to be:
        """
        CommonName: group01
        GidNumber: \"\d+\"
        MemberUid: uix00002
        ObjectClass:
        - top
        - posixGroup
        """
        And I expect the exit code to be 0

    Scenario: I submit a search request to find a specific user by id
        Given I perform bind with following parameters:
        | --user     | cn=admin,dc=mock,dc=ad,dc=com |
        | --password | admin                         |
        | --url      | ldap://localhost:389          |
        And I plan to execute "get" command with following parameters:
        | --path   | dc=mock,dc=ad,dc=com |
        | --select | *                    |
        And I plan to execute "user" command with following parameters:
        | --user-id | uix00001 |
        When I execute the application
        Then I expect the output to be:
        """
        CommonName: uix00001
        Gecos: uix00001
        GidNumber: \"100\"
        HomeDirectory: /home/uix00001
        LoginShell: /bin/bash
        ObjectClass:
        - top
        - account
        - posixAccount
        - shadowAccount
        ShadowLastChange: \"0\"
        ShadowMax: \"0\"
        ShadowWarning: \"0\"
        Uid: uix00001
        UidNumber: \"\d+\"
        UserPassword: new-password
        """
        And I expect the exit code to be 0

    Scenario: I submit a search request to find a specific group by id
        Given I perform bind with following parameters:
        | --user     | cn=admin,dc=mock,dc=ad,dc=com |
        | --password | admin                         |
        | --url      | ldap://localhost:389          |
        And I plan to execute "get" command with following parameters:
        | --path   | dc=mock,dc=ad,dc=com |
        | --select | *                    |
        And I plan to execute "group" command with following parameters:
        | --group-id | group01 |
        When I execute the application
        Then I expect the output to be:
        """
        CommonName: group01
        GidNumber: \"\d+\"
        MemberUid: uix00002
        ObjectClass:
        - top
        - posixGroup
        """
        And I expect the exit code to be 0
