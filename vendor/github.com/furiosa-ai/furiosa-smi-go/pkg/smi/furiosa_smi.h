#ifndef __furiosa_smi_h__
#define __furiosa_smi_h__

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

#define FURIOSA_SMI_MAX_PATH_SIZE 256

#define FURIOSA_SMI_MAX_DEVICE_FILE_SIZE 64

#define FURIOSA_SMI_MAX_CORE_STATUS_SIZE 128

#define FURIOSA_SMI_MAX_PE_SIZE 64

#define FURIOSA_SMI_MAX_DEVICE_HANDLE_SIZE 64

#define FURIOSA_SMI_MAX_CSTR_SIZE 96

/// \brief Represent an architecture type of device
typedef enum {
  /// Warboy architecture
  FURIOSA_SMI_ARCH_WARBOY = 0,
  /// RNGD architecture
  FURIOSA_SMI_ARCH_RNGD,
  /// RNGD-Max architecture
  FURIOSA_SMI_ARCH_RNGD_MAX,
  /// RNGD-S architecture
  FURIOSA_SMI_ARCH_RNGD_S,
} FuriosaSmiArch;

/// \brief Represent a core status
typedef enum {
  /// Device is available
  FURIOSA_SMI_CORE_STATUS_AVAILABLE = 0,
  /// Device is occupied
  FURIOSA_SMI_CORE_STATUS_OCCUPIED,
} FuriosaSmiCoreStatus;

/// \brief Represent a device link type
typedef enum {
  /// Unknown link type
  FURIOSA_SMI_DEVICE_TO_DEVICE_LINK_TYPE_UNKNOWN = 0,
  /// Link type under same machine
  FURIOSA_SMI_DEVICE_TO_DEVICE_LINK_TYPE_INTERCONNECT = 10,
  /// Link type under same cpu
  FURIOSA_SMI_DEVICE_TO_DEVICE_LINK_TYPE_CPU = 20,
  /// Link type under same switch
  FURIOSA_SMI_DEVICE_TO_DEVICE_LINK_TYPE_BRIDGE = 30,
  /// Link type under same socket
  FURIOSA_SMI_DEVICE_TO_DEVICE_LINK_TYPE_NOC = 70,
} FuriosaSmiDeviceToDeviceLinkType;

/// \brief Represent a governor profile
typedef enum {
  /// 'performance' governor profile
  FURIOSA_SMI_GOVERNOR_PROFILE_PERFORMANCE = 0,
  /// 'powersave' governor profile
  FURIOSA_SMI_GOVERNOR_PROFILE_POWER_SAVE,
} FuriosaSmiGovernorProfile;

/// \brief Represent a return status
typedef enum {
  /// When a function call is successful.
  FURIOSA_SMI_RETURN_CODE_OK = 0,
  /// When an invalid argument is given.
  FURIOSA_SMI_RETURN_CODE_INVALID_ARGUMENT_ERROR,
  /// When a null pointer is given to output buffer.
  FURIOSA_SMI_RETURN_CODE_NULL_POINTER_ERROR,
  /// When a data exceeds the maximum buffer size.
  FURIOSA_SMI_RETURN_CODE_MAX_BUFFER_SIZE_EXCEED_ERROR,
  /// When a device is not found with the given option.
  FURIOSA_SMI_RETURN_CODE_DEVICE_NOT_FOUND_ERROR,
  /// When a device is lost.
  FURIOSA_SMI_RETURN_CODE_DEVICE_LOST_ERROR,
  /// When a device state is busy.
  FURIOSA_SMI_RETURN_CODE_DEVICE_BUSY_ERROR,
  /// When a certain operation is failed by an unexpected io error.
  FURIOSA_SMI_RETURN_CODE_IO_ERROR,
  /// When a certain operation is failed by a permission deny.
  FURIOSA_SMI_RETURN_CODE_PERMISSION_DENIED_ERROR,
  /// When an arch is unknown.
  FURIOSA_SMI_RETURN_CODE_UNKNOWN_ARCH_ERROR,
  /// When a driver is incompatible.
  FURIOSA_SMI_RETURN_CODE_INCOMPATIBLE_DRIVER_ERROR,
  /// When a retrieved value is invalid.
  FURIOSA_SMI_RETURN_CODE_UNEXPECTED_VALUE_ERROR,
  /// When a certain parsing operation is failed.
  FURIOSA_SMI_RETURN_CODE_PARSE_ERROR,
  /// When a reason is unknown.
  FURIOSA_SMI_RETURN_CODE_UNKNOWN_ERROR,
  /// When an internal operation is failed.
  FURIOSA_SMI_RETURN_CODE_INTERNAL_ERROR,
  /// When the system is not initialized.
  FURIOSA_SMI_RETURN_CODE_UNINITIALIZED_ERROR,
  /// When a context cannot be captured.
  FURIOSA_SMI_RETURN_CODE_CONTEXT_ERROR,
  /// When a certain operation is not supported.
  FURIOSA_SMI_RETURN_CODE_NOT_SUPPORTED_ERROR,
} FuriosaSmiReturnCode;

typedef uint32_t FuriosaSmiDeviceHandle;

/// \brief Represent a device handle list in the system.
typedef struct {
  uint32_t count;
  FuriosaSmiDeviceHandle device_handles[FURIOSA_SMI_MAX_DEVICE_HANDLE_SIZE];
} FuriosaSmiDeviceHandles;

/// \brief Represent the version of device component
typedef struct {
  uint32_t major;
  uint32_t minor;
  uint32_t patch;
  char metadata[FURIOSA_SMI_MAX_CSTR_SIZE];
} FuriosaSmiVersion;

/// \brief Represent a device information
typedef struct {
  uint32_t index;
  FuriosaSmiArch arch;
  uint32_t core_num;
  uint32_t numa_node;
  char name[FURIOSA_SMI_MAX_CSTR_SIZE];
  char serial[FURIOSA_SMI_MAX_CSTR_SIZE];
  char uuid[FURIOSA_SMI_MAX_CSTR_SIZE];
  char bdf[FURIOSA_SMI_MAX_CSTR_SIZE];
  uint16_t major;
  uint16_t minor;
  FuriosaSmiVersion firmware_version;
  FuriosaSmiVersion pert_version;
} FuriosaSmiDeviceInfo;

/// \brief Represent a device file information.
typedef struct {
  uint32_t core_start;
  uint32_t core_end;
  char path[FURIOSA_SMI_MAX_PATH_SIZE];
} FuriosaSmiDeviceFile;

/// \brief Represent a device file list of device.
typedef struct {
  uint32_t count;
  FuriosaSmiDeviceFile device_files[FURIOSA_SMI_MAX_DEVICE_HANDLE_SIZE];
} FuriosaSmiDeviceFiles;

typedef struct {
  uint32_t core;
  FuriosaSmiCoreStatus status;
} FuriosaSmiPeStatus;

/// \brief Represent a core status list of device.
typedef struct {
  uint32_t count;
  FuriosaSmiPeStatus core_status[FURIOSA_SMI_MAX_CORE_STATUS_SIZE];
} FuriosaSmiCoreStatuses;

typedef struct {
  uint32_t core;
  uint32_t frequency;
} FuriosaSmiPeFrequency;

/// \brief Represent a core frequency information (MHz).
typedef struct {
  uint32_t pe_count;
  FuriosaSmiPeFrequency pe[FURIOSA_SMI_MAX_PE_SIZE];
} FuriosaSmiCoreFrequency;

/// \brief Represent a memory frequency information (MHz).
typedef struct {
  uint32_t frequency;
} FuriosaSmiMemoryFrequency;

/// \brief Represent a PE performance counter.
typedef struct {
  long timestamp;
  uint32_t core;
  uint64_t cycle_count;
  uint64_t task_execution_cycle;
} FuriosaSmiPePerformanceCounter;

/// \brief Represent a device performance counter
typedef struct {
  uint32_t pe_count;
  FuriosaSmiPePerformanceCounter pe_performance_counters[FURIOSA_SMI_MAX_PE_SIZE];
} FuriosaSmiDevicePerformanceCounter;

/// \brief Represent a power consumption of device.
typedef struct {
  double rms_total;
} FuriosaSmiDevicePowerConsumption;

/// \brief Represent a temperature information of device.
typedef struct {
  double soc_peak;
  double ambient;
} FuriosaSmiDeviceTemperature;

typedef char FuriosaSmiBdf[FURIOSA_SMI_MAX_CSTR_SIZE];

typedef struct {
  uint32_t count;
  FuriosaSmiBdf bdfs[FURIOSA_SMI_MAX_DEVICE_HANDLE_SIZE];
} FuriosaSmiDisabledDevices;

/// \brief Represent a PCIe device information.
typedef struct {
  uint16_t device_id;
  uint16_t subsystem_vendor_id;
  uint16_t subsystem_device_id;
  uint8_t revision_id;
  uint8_t class_id;
  uint8_t sub_class_id;
} FuriosaSmiPcieDeviceInfo;

/// \brief Represent a PCIe link information.
typedef struct {
  uint8_t pcie_gen_status;
  uint32_t link_width_status;
  double link_speed_status;
  uint32_t max_link_width_capability;
  double max_link_speed_capability;
} FuriosaSmiPcieLinkInfo;

/// \brief Represent a SR-IOV information.
typedef struct {
  uint32_t sriov_total_vfs;
  uint32_t sriov_enabled_vfs;
} FuriosaSmiSriovInfo;

/// \brief Represent a PCIe root complex information.
typedef struct {
  uint16_t domain;
  uint8_t bus;
} FuriosaSmiPcieRootComplexInfo;

/// \brief Represent a PCIe switch information. If switch doesn't exist, all values are max value.
typedef struct {
  uint16_t domain;
  uint8_t bus;
  uint8_t device;
  uint8_t function;
} FuriosaSmiPcieSwitchInfo;

/// @defgroup Initialize Initialize
/// @brief Initialize module for Furiosa smi.
/// @{

/// \brief Initialize Furiosa SMI library.
///
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_init(void);

/// @}

/// @defgroup Device Device
/// @brief Device module for Furiosa smi.
/// @{

/// @brief Get all device handles of Furiosa NPU devices in the system.
/// @param[out] out_handles output buffer for pointer to FuriosaSmiDeviceHandles.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_device_handles(FuriosaSmiDeviceHandles *out_handles);

/// \brief Get a device handle of Furiosa NPU device by uuid.
///
/// @param uuid uuid of Furiosa NPU device.
/// @param[out] out_handle output buffer for pointer to FuriosaSmiDeviceHandle of given uuid.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_device_handle_by_uuid(const char *uuid,
                                                           FuriosaSmiDeviceHandle *out_handle);

/// \brief Get a device handle of Furiosa NPU device by serial.
///
/// @param serial serial of Furiosa NPU device.
/// @param[out] out_handle output buffer for pointer to FuriosaSmiDeviceHandle of given serial.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_device_handle_by_serial(const char *serial,
                                                             FuriosaSmiDeviceHandle *out_handle);

/// \brief Get a device handle of Furiosa NPU device by bdf.
///
/// @param bdf bdf of Furiosa NPU device.
/// @param[out] out_handle output buffer for pointer to FuriosaSmiDeviceHandle of given bdf.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_device_handle_by_bdf(const char *bdf,
                                                          FuriosaSmiDeviceHandle *out_handle);

/// \brief Get a device information of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_device_info output buffer for pointer to FuriosaSmiDeviceInfo.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_device_info(FuriosaSmiDeviceHandle handle,
                                                 FuriosaSmiDeviceInfo *out_device_info);

/// \brief Get a device file list of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_device_files output buffer for pointer to FuriosaSmiDeviceFiles.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_device_files(FuriosaSmiDeviceHandle handle,
                                                  FuriosaSmiDeviceFiles *out_device_files);

/// \brief Get a core status list of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_core_status output buffer for pointer to FuriosaSmiCoreStatuses.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_device_core_status(FuriosaSmiDeviceHandle handle,
                                                        FuriosaSmiCoreStatuses *out_core_status);

/// \brief Get a liveness of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_liveness output buffer for pointer to boolean representing the liveness of device.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_device_liveness(FuriosaSmiDeviceHandle handle,
                                                     bool *out_liveness);

/// @}

/// @defgroup Topology Topology
/// @brief Topology module for Furiosa smi.
/// @{

/// \brief Get a device link type between two Furiosa NPU devices.
///
/// @param handle1 handle of Furiosa NPU device 1.
/// @param handle2 handle of Furiosa NPU device 2.
/// @param[out] out_link_type output buffer for pointer to FuriosaSmiDeviceToDeviceLinkType.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_device_to_device_link_type(FuriosaSmiDeviceHandle handle1,
                                                                FuriosaSmiDeviceHandle handle2,
                                                                FuriosaSmiDeviceToDeviceLinkType *out_link_type);

/// \brief Checks if two Furiosa NPU devices are P2P accessible.
///
/// @param handle1 handle of Furiosa NPU device 1.
/// @param handle2 handle of Furiosa NPU device 2.
/// @param[out] out_accessible output buffer for pointer to boolean result.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_p2p_accessible(FuriosaSmiDeviceHandle handle1,
                                                    FuriosaSmiDeviceHandle handle2,
                                                    bool *out_accessible);

/// @}

/// @defgroup System System
/// @brief System module for Furiosa smi.
/// @{

/// \brief Get a driver information of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_driver_info output buffer for pointer to FuriosaSmiVersion.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_driver_info(FuriosaSmiVersion *out_driver_info);

/// \brief Get a core frequency of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_core_frequency_info output buffer for pointer to FuriosaSmiCoreFrequency.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_core_frequency(FuriosaSmiDeviceHandle handle,
                                                    FuriosaSmiCoreFrequency *out_core_frequency_info);

/// \brief Get a memory frequency of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_memory_frequency_info output buffer for pointer to FuriosaSmiMemoryFrequency.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_memory_frequency(FuriosaSmiDeviceHandle handle,
                                                      FuriosaSmiMemoryFrequency *out_memory_frequency_info);

/// \brief Get a performance counter of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param core core index of Furiosa NPU device.
/// @param[out] out_performance_counter_info output buffer for pointer to FuriosaSmiPePerformanceCounter.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_device_performance_counter(FuriosaSmiDeviceHandle handle,
                                                                FuriosaSmiDevicePerformanceCounter *out_performance_counter_info);

/// \brief Get a power consumption of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_power_consumption output buffer for pointer to FuriosaSmiDevicePowerConsumption.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_device_power_consumption(FuriosaSmiDeviceHandle handle,
                                                              FuriosaSmiDevicePowerConsumption *out_power_consumption);

/// \brief Get a power consumption of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_power_consumption output buffer for pointer to FuriosaSmiDevicePowerConsumption.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_device_temperature(FuriosaSmiDeviceHandle handle,
                                                        FuriosaSmiDeviceTemperature *out_temperature);

/// \brief Get a governor state of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_governor_profile output buffer for pointer to FuriosaSmiGovernorProfile.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_governor_profile(FuriosaSmiDeviceHandle handle,
                                                      FuriosaSmiGovernorProfile *out_governor_profile);

/// \brief Set a governor state into Furiosa NPU device. This requires root privileges.
///
/// @param handle handle of Furiosa NPU device.
/// @param[in] governor_profile input buffer for pointer to FuriosaSmiGovernorProfile.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_set_governor_profile(FuriosaSmiDeviceHandle handle,
                                                      FuriosaSmiGovernorProfile governor_profile);

/// \brief Bind a Furiosa NPU device. This requires root privileges.
///
/// @param handle handle of Furiosa NPU device.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_enable_device(FuriosaSmiDeviceHandle handle);

/// \brief Unbind a Furiosa NPU device. This requires root privileges.
///
/// @param handle handle of Furiosa NPU device.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_disable_device(FuriosaSmiDeviceHandle handle);

/// \brief Get a list of disabled Furiosa NPU devices.
///
/// @param[out] out_disabled_devices output buffer for pointer to FuriosaSmiDisabledDevices.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_disabled_devices(FuriosaSmiDisabledDevices *out_disabled_devices);

/// \brief Get a PCIe information of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_pcie_device_info output buffer for pointer to FuriosaSmiPcieDeviceInfo.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_pcie_device_info(FuriosaSmiDeviceHandle handle,
                                                      FuriosaSmiPcieDeviceInfo *out_pcie_device_info);

/// \brief Get a PCIe link information of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_pcie_link_info output buffer for pointer to FuriosaSmiPcieLinkInfo.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_pcie_link_info(FuriosaSmiDeviceHandle handle,
                                                    FuriosaSmiPcieLinkInfo *out_pcie_link_info);

/// \brief Get a SR-IOV information of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_sriov_info output buffer for pointer to FuriosaSmiSriovInfo.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_sriov_info(FuriosaSmiDeviceHandle handle,
                                                FuriosaSmiSriovInfo *out_sriov_info);

/// \brief Get a PCIe root complex information of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_root_complex_info output buffer for pointer to FuriosaSmiPcieRootComplexInfo.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_pcie_root_complex_info(FuriosaSmiDeviceHandle handle,
                                                            FuriosaSmiPcieRootComplexInfo *out_root_complex_info);

/// \brief Get a PCIe switch information of Furiosa NPU device.
///
/// @param handle handle of Furiosa NPU device.
/// @param[out] out_pcie_switch_info output buffer for pointer to FuriosaSmiPcieSwitchInfo.
/// @return FURIOSA_SMI_RETURN_CODE_OK if successful, see `FuriosaSmiReturnCode` for error cases.
FuriosaSmiReturnCode furiosa_smi_get_pcie_switch_info(FuriosaSmiDeviceHandle handle,
                                                      FuriosaSmiPcieSwitchInfo *out_pcie_switch_info);

/// @}



#endif /* __furiosa_smi_h__ */
