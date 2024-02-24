Feature: App version

    Scenario: I check version of the app
        Given I plan to execute "version" command with following parameters:
        |||
        When I execute the application
        Then I expect the output to be:
        """
        Version: v\d+.\d+.\d+ .*
        Built at: .+
        Executable path: .+
        """
        And I expect the exit code to be 0