.. _implementedservice_appconfig:

.. |start-h3| raw:: html

    <h3>

.. |end-h3| raw:: html

    </h3>

=========
appconfig
=========

.. autoclass:: moto.appconfig.models.AppConfigBackend

|start-h3| Example usage |end-h3|

.. sourcecode:: python

            @mock_appconfig
            def test_appconfig_behaviour():
                client = boto3.client("appconfig", region_name="us-east-1")
                ...



|start-h3| Implemented features for this service |end-h3|

- [X] create_application
- [X] create_configuration_profile
- [X] create_hosted_configuration_version
- [X] get_hosted_configuration_version
