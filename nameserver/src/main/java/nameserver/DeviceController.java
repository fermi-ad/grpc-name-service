package nameserver;

import com.google.protobuf.Empty;
import java.util.Arrays;
import java.util.ArrayList;
import io.quarkus.hibernate.orm.panache.PanacheQuery;
import io.quarkus.panache.common.Parameters;
import jakarta.transaction.Transactional;
import nameserver.model.Device;
import java.util.List;

/*Device API implementation */
public class DeviceController {
    
    public static Device getDeviceByName(String name) {        
        Device device = Device.find("name", name).firstResult();
        if (device == null) {
            throw new RuntimeException("Device with name " + name + " not found");
        }
        return device;
    }
    
    static proto.Device.Builder createDeviceBuilder(Device device) {
        var deviceBuilder = proto.Device.newBuilder();
        
        deviceBuilder.setName(device.getName())
                .setDescription(device.getDescription())
                .setNodeHostname(device.getNode().getHostname());
        
        return deviceBuilder;
    }    
    @Transactional
    public static proto.DeviceResponse createDevice(proto.CreateDeviceRequest request) {        
        var reqDevice = request.getDevice();
        
        var node = NodeController.getNodeByHostname(reqDevice.getNodeHostname());
        Device device = new Device(reqDevice.getName(),
                             reqDevice.getDescription(),
                             node);
        device.persist();

        var resDeviceBuilder = createDeviceBuilder(device);        
        return proto.DeviceResponse.newBuilder().setDevice(resDeviceBuilder).build();
    }

    @Transactional
    public static proto.DeviceResponse updateDevice(proto.UpdateDeviceRequest request) {
        var reqDevice = request.getDevice();
        Device device = getDeviceByName(reqDevice.getName());        
        if (reqDevice.getNodeHostname() != "") {
            device.setNode(NodeController.getNodeByHostname(reqDevice.getNodeHostname()));
        }
        if(reqDevice.getDescription() != "") {
            device.setDescription(reqDevice.getDescription());
        }
        device.persist();
        var resDeviceBuilder = createDeviceBuilder(device);        

        return proto.DeviceResponse.newBuilder().setDevice(resDeviceBuilder).build();                
    }

    public static proto.DeviceListResponse listDevices(proto.ListDevicesRequest request) {
        var queryParams = new ArrayList<List<String>>();
        if (!request.getName().isEmpty()) {
            queryParams.add(Arrays.asList("name",request.getName()));
        }
        if (!request.getNodeHostname().isEmpty()) {
            queryParams.add(Arrays.asList("node.hostname",request.getNodeHostname()));
        }
        Parameters parameters = new Parameters();   
            
        var queryString = ControllerUtil.createFindParam(queryParams, parameters);        
        PanacheQuery<Device> query;
        query = Device.find(queryString, parameters);
        
        var pag = ControllerUtil.readPagination(request.getPagination());                
        List<Device> devices = query.page(pag[0] - 1, pag[1]).list();        
        proto.DeviceListResponse.Builder response = proto.DeviceListResponse.newBuilder();
        for (Device device : devices) {
            var devBuilder = createDeviceBuilder(device);            
            response.addDevices(devBuilder.build());
        }
        response.setPagination(ControllerUtil.createPaginationResponse(pag, query.count()));
        return response.build();
    }

    @Transactional
    public static Empty deleteDevice(proto.DeleteDeviceRequest request) {  
        Device.delete("name", request.getName());
        return Empty.newBuilder().build();
    }     

    public static proto.DeviceResponse getDevice(proto.GetDeviceRequest request) {
        Device device = getDeviceByName(request.getName());
        var deviceBuilder = createDeviceBuilder(device);
        return proto.DeviceResponse.newBuilder().setDevice(deviceBuilder).build();
    }
}

