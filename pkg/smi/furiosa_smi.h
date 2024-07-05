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

#define FURIOSA_SMI_MAX_DRIVER_INFO_SIZE 24

#define FURIOSA_SMI_MAX_DEVICE_HANDLE_SIZE 64

#define FURIOSA_SMI_MAX_CSTR_SIZE 96

typedef enum {
  FURIOSA_SMI_ARCH_WARBOY = 0,
  FURIOSA_SMI_ARCH_RNGD,
  FURIOSA_SMI_ARCH_RNGD_MAX,
  FURIOSA_SMI_ARCH_RNGD_S,
} FuriosaSmiArch;

typedef enum {
  FURIOSA_SMI_CORE_STATUS_AVAILABLE = 0,
  FURIOSA_SMI_CORE_STATUS_OCCUPIED,
} FuriosaSmiCoreStatus;

typedef enum {
  FURIOSA_SMI_DEVICE_TO_DEVICE_LINK_TYPE_UNKNOWN = 0,
  FURIOSA_SMI_DEVICE_TO_DEVICE_LINK_TYPE_INTERCONNECT = 10,
  FURIOSA_SMI_DEVICE_TO_DEVICE_LINK_TYPE_CPU = 20,
  FURIOSA_SMI_DEVICE_TO_DEVICE_LINK_TYPE_BRIDGE = 30,
  FURIOSA_SMI_DEVICE_TO_DEVICE_LINK_TYPE_NOC = 70,
} FuriosaSmiDeviceToDeviceLinkType;

typedef enum {
  FURIOSA_SMI_RETURN_CODE_OK = 0,
  FURIOSA_SMI_RETURN_CODE_INVALID_ARGUMENT_ERROR,
  FURIOSA_SMI_RETURN_CODE_NULL_POINTER_ERROR,
  FURIOSA_SMI_RETURN_CODE_MAX_BUFFER_SIZE_EXCEED_ERROR,
  FURIOSA_SMI_RETURN_CODE_DEVICE_NOT_FOUND_ERROR,
  FURIOSA_SMI_RETURN_CODE_DEVICE_BUSY_ERROR,
  FURIOSA_SMI_RETURN_CODE_IO_ERROR,
  FURIOSA_SMI_RETURN_CODE_PERMISSION_DENIED_ERROR,
  FURIOSA_SMI_RETURN_CODE_UNKNOWN_ARCH_ERROR,
  FURIOSA_SMI_RETURN_CODE_INCOMPATIBLE_DRIVER_ERROR,
  FURIOSA_SMI_RETURN_CODE_UNEXPECTED_VALUE_ERROR,
  FURIOSA_SMI_RETURN_CODE_PARSE_ERROR,
  FURIOSA_SMI_RETURN_CODE_UNKNOWN_ERROR,
  FURIOSA_SMI_RETURN_CODE_INTERNAL_ERROR,
  FURIOSA_SMI_RETURN_CODE_UNINITIALIZED_ERROR,
  FURIOSA_SMI_RETURN_CODE_CONTEXT_ERROR,
} FuriosaSmiReturnCode;

typedef struct FuriosaSmiObserver FuriosaSmiObserver;

typedef uint32_t FuriosaSmiDeviceHandle;

typedef struct {
  uint32_t count;
  FuriosaSmiDeviceHandle device_handles[FURIOSA_SMI_MAX_DEVICE_HANDLE_SIZE];
} FuriosaSmiDeviceHandles;

typedef FuriosaSmiObserver* FuriosaSmiObserverInstance;

typedef struct {
  FuriosaSmiArch arch;
  uint32_t major;
  uint32_t minor;
  uint32_t patch;
  char metadata[FURIOSA_SMI_MAX_CSTR_SIZE];
} FuriosaSmiVersion;

typedef struct {
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
  FuriosaSmiVersion driver_version;
} FuriosaSmiDeviceInfo;

typedef struct {
  uint32_t core_start;
  uint32_t core_end;
  char path[FURIOSA_SMI_MAX_PATH_SIZE];
} FuriosaSmiDeviceFile;

typedef struct {
  uint32_t count;
  FuriosaSmiDeviceFile device_files[FURIOSA_SMI_MAX_DEVICE_HANDLE_SIZE];
} FuriosaSmiDeviceFiles;

typedef struct {
  uint32_t count;
  FuriosaSmiCoreStatus core_status[FURIOSA_SMI_MAX_CORE_STATUS_SIZE];
} FuriosaSmiCoreStatuses;

typedef struct {
  uint32_t axi_post_error_count;
  uint32_t axi_fetch_error_count;
  uint32_t axi_discard_error_count;
  uint32_t axi_doorbell_error_count;
  uint32_t pcie_post_error_count;
  uint32_t pcie_fetch_error_count;
  uint32_t pcie_discard_error_count;
  uint32_t pcie_doorbell_error_count;
  uint32_t device_error_count;
} FuriosaSmiDeviceErrorInfo;

typedef struct {
  uint32_t count;
  FuriosaSmiVersion driver_info[FURIOSA_SMI_MAX_DRIVER_INFO_SIZE];
} FuriosaSmiDriverInfo;

typedef struct {
  uint32_t core_count;
  uint32_t cores[FURIOSA_SMI_MAX_PE_SIZE];
  uint32_t time_window_mil;
  double pe_usage_percentage;
} FuriosaSmiPeUtilization;

typedef struct {
  uint64_t total_bytes;
  uint64_t in_use_bytes;
} FuriosaSmiMemoryUtilization;

typedef struct {
  uint32_t pe_count;
  FuriosaSmiPeUtilization pe[FURIOSA_SMI_MAX_PE_SIZE];
  FuriosaSmiMemoryUtilization memory;
} FuriosaSmiDeviceUtilization;

typedef struct {
  double rms_total;
} FuriosaSmiDevicePowerConsumption;

typedef struct {
  double soc_peak;
  double ambient;
} FuriosaSmiDeviceTemperature;

FuriosaSmiReturnCode furiosa_smi_get_device_handles(FuriosaSmiDeviceHandles *out_handles);

FuriosaSmiReturnCode furiosa_smi_create_observer(FuriosaSmiObserverInstance *out_observer_instance);

FuriosaSmiReturnCode furiosa_smi_destroy_observer(FuriosaSmiObserverInstance *p_observer_instance);

FuriosaSmiReturnCode furiosa_smi_get_device_handle_by_uuid(const char *uuid,
                                                           FuriosaSmiDeviceHandle *out_handle);

FuriosaSmiReturnCode furiosa_smi_get_device_handle_by_serial(const char *serial,
                                                             FuriosaSmiDeviceHandle *out_handle);

FuriosaSmiReturnCode furiosa_smi_get_device_handle_by_bdf(const char *bdf,
                                                          FuriosaSmiDeviceHandle *out_handle);

FuriosaSmiReturnCode furiosa_smi_get_device_info(FuriosaSmiDeviceHandle handle,
                                                 FuriosaSmiDeviceInfo *out_device_info);

FuriosaSmiReturnCode furiosa_smi_get_device_files(FuriosaSmiDeviceHandle handle,
                                                  FuriosaSmiDeviceFiles *out_device_files);

FuriosaSmiReturnCode furiosa_smi_get_device_core_status(FuriosaSmiDeviceHandle handle,
                                                        FuriosaSmiCoreStatuses *out_core_status);

FuriosaSmiReturnCode furiosa_smi_get_device_liveness(FuriosaSmiDeviceHandle handle,
                                                     bool *out_liveness);

FuriosaSmiReturnCode furiosa_smi_get_device_error_info(FuriosaSmiDeviceHandle handle,
                                                       FuriosaSmiDeviceErrorInfo *out_error_info);

FuriosaSmiReturnCode furiosa_smi_get_driver_info(FuriosaSmiDriverInfo *out_driver_info);

FuriosaSmiReturnCode furiosa_smi_get_device_utilization(FuriosaSmiObserverInstance observer_instance,
                                                        FuriosaSmiDeviceHandle handle,
                                                        FuriosaSmiDeviceUtilization *out_utilization_info);

FuriosaSmiReturnCode furiosa_smi_get_device_power_consumption(FuriosaSmiDeviceHandle handle,
                                                              FuriosaSmiDevicePowerConsumption *out_power_consumption);

FuriosaSmiReturnCode furiosa_smi_get_device_temperature(FuriosaSmiDeviceHandle handle,
                                                        FuriosaSmiDeviceTemperature *out_temperature);

FuriosaSmiReturnCode furiosa_smi_get_device_to_device_link_type(FuriosaSmiDeviceHandle handle1,
                                                                FuriosaSmiDeviceHandle handle2,
                                                                FuriosaSmiDeviceToDeviceLinkType *out_link_type);

#endif /* __furiosa_smi_h__ */
