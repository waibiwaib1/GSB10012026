describe('AppleSimUtils', () => {
  let exec;
  let AppleSimUtils;
  let appleSimUtils;

  beforeEach(() => {
    jest.mock('../../utils/exec');
    jest.mock('../../utils/environment');
    exec = require('../../utils/exec');
    AppleSimUtils = require('./AppleSimUtils');
    appleSimUtils = new AppleSimUtils();
  });

  const simulatorList = ({ deviceTypes = [deviceType()], runtimes = [runtime()] } = {}) => ({
    stdout: JSON.stringify({
      devicetypes: deviceTypes,
      runtimes,
    }),
  });

  const deviceType = () => ({
    name: 'iPhone X',
    identifier: 'com.apple.CoreSimulator.SimDeviceType.iPhone-X',
  });

  const runtime = (runtimeProps = {}) => ({
    name: 'iOS 12.2',
    identifier: 'com.apple.CoreSimulator.SimRuntime.iOS-12-2',
    version: '12.2',
    availability: '(available)',
    ...runtimeProps,
  });

  it('creates a simulator using the newest available iOS runtime', async () => {
    exec.execWithRetriesAndLogs
      .mockResolvedValueOnce(simulatorList({
        runtimes: [
          runtime({ name: 'tvOS 13.0', identifier: 'tv-runtime', version: '13.0' }),
          runtime({ name: 'iOS 12.2', version: '12.2' }),
        ],
      }))
      .mockResolvedValueOnce({ stdout: 'new-device-udid\n' });

    const udid = await appleSimUtils.create('iPhone X');

    expect(exec.execWithRetriesAndLogs).toHaveBeenNthCalledWith(
      2,
      '/usr/bin/xcrun simctl create "iPhone X-Detox" "com.apple.CoreSimulator.SimDeviceType.iPhone-X" "com.apple.CoreSimulator.SimRuntime.iOS-12-2"',
      { verbosity: 'normal' },
      {},
      1
    );
    expect(udid).toBe('new-device-udid');
  });

  it('creates an OS-specific simulator', async () => {
    exec.execWithRetriesAndLogs
      .mockResolvedValueOnce(simulatorList({
        runtimes: [
          runtime({ name: 'iOS 11.4', identifier: 'ios-11-4', version: '11.4' }),
          runtime(),
        ],
      }))
      .mockResolvedValueOnce({ stdout: 'new-device-udid\n' });

    await appleSimUtils.create('iPhone X, iOS 11.4');

    expect(exec.execWithRetriesAndLogs).toHaveBeenNthCalledWith(
      2,
      expect.stringContaining('"ios-11-4"'),
      expect.anything(),
      expect.anything(),
      expect.anything()
    );
  });

  it('fails with the supported device types when the type is missing', async () => {
    exec.execWithRetriesAndLogs.mockResolvedValueOnce(simulatorList({ deviceTypes: [] }));

    await expect(appleSimUtils.create('iPhone XR'))
      .rejects
      .toThrow(`Can't find a simulator device type named "iPhone XR", run 'xcrun simctl list devicetypes' to list supported device types.`);
  });

  it('fails with the supported runtimes when the runtime is missing', async () => {
    exec.execWithRetriesAndLogs.mockResolvedValueOnce(simulatorList({ runtimes: [] }));

    await expect(appleSimUtils.create('iPhone X'))
      .rejects
      .toThrow(`Can't find an available iOS runtime, run 'xcrun simctl list runtimes' to list supported runtimes.`);
  });
});
