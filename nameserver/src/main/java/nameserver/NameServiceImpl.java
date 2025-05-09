package nameserver;
import com.google.protobuf.Empty;

import io.quarkus.grpc.GrpcService;
import io.smallrye.common.annotation.Blocking;
import io.smallrye.mutiny.Multi;
import io.smallrye.mutiny.Uni;
import proto.AddChannelAccessControlRequest;
import proto.AddChannelAlarmRequest;
import proto.AddChannelTransformRequest;
import proto.AlarmTypeListResponse;
import proto.AlarmTypeResponse;
import proto.ChannelAccessControlResponse;
import proto.ChannelAlarmResponse;
import proto.ChannelResponse;
import proto.ChannelTransformResponse;
import proto.CreateAlarmTypeRequest;
import proto.CreateRoleRequest;
import proto.DeleteAlarmTypeRequest;
import proto.DeleteChannelAccessControlRequest;
import proto.DeleteChannelAlarmRequest;
import proto.DeleteChannelTransformRequest;
import proto.DeleteRoleRequest;
import proto.DeviceResponse;
import proto.GetChannelRequest;
import proto.GetDeviceRequest;
import proto.GetLocationRequest;
import proto.GetNodeRequest;
import proto.LocationResponse;
import proto.NodeResponse;
import proto.RoleListResponse;
import proto.RoleResponse;
import proto.UpdateAlarmTypeRequest;
import proto.UpdateChannelAccessControlRequest;
import proto.UpdateChannelAlarmRequest;
import proto.UpdateChannelTransformRequest;
import proto.UpdateRoleRequest;
import io.grpc.StatusRuntimeException;
import io.grpc.Status;
import java.util.logging.Logger;
import jakarta.annotation.security.RolesAllowed;;

@GrpcService
public class NameServiceImpl implements proto.NameService {

    static final Logger logger = java.util.logging.Logger.getLogger(NameServiceImpl.class.getName());

    private Throwable getRootCause(Throwable throwable) {
        Throwable cause = throwable;
        while (cause.getCause() != null) {
            cause = cause.getCause();
        }
        return cause;
    }
    private void printCallStack(Throwable throwable) {
        int limit = 10;
        int i = 0;
        // Log the throwable message
        logger.severe("Exception occurred: " + throwable.getMessage());
    
        // Print the stack trace to the logger
        for (StackTraceElement element : throwable.getStackTrace()) {
            if (i >= limit) {
                break; // Limit the number of stack trace elements to log
            }   
            logger.severe("\tat " + element.toString());
            i++;
        }
    }
    private <T> Uni<T> createUni(java.util.function.Supplier<T> action) {
        return Uni.createFrom().item(action)
        .onFailure().recoverWithItem((Throwable throwable) -> {
            Throwable rootCause = getRootCause(throwable);

            printCallStack(throwable);
            if (rootCause instanceof StatusRuntimeException) {                
                throw Status.fromThrowable(rootCause).asRuntimeException();
            } else if (rootCause instanceof org.hibernate.exception.ConstraintViolationException) {
                throw Status.ALREADY_EXISTS
                    .withDescription(rootCause.getMessage())
                    .asRuntimeException();
            } else if (rootCause instanceof org.hibernate.exception.JDBCConnectionException) {
                throw Status.UNAVAILABLE
                    .withDescription(rootCause.getMessage())
                    .asRuntimeException();
            } else if (rootCause instanceof org.hibernate.exception.SQLGrammarException) {
                throw Status.INVALID_ARGUMENT
                    .withDescription(rootCause.getMessage())
                    .asRuntimeException();
            } else {
                throw Status.INTERNAL
                    .withDescription(rootCause.getMessage())
                    .asRuntimeException();
            }
        });
    }
    @RolesAllowed("admin")
    @Override
    @Blocking
    public Uni<proto.LocationTypeResponse> createLocationType(proto.CreateLocationTypeRequest request) {       
        return createUni(() -> LocationTypeController.createLocationType(request));
    }

    @RolesAllowed("admin")    
    @Override
    @Blocking
    public Uni<proto.LocationTypeResponse> updateLocationType(proto.UpdateLocationTypeRequest request) {
        return createUni(() -> LocationTypeController.updateLocationType(request));
    }

    @RolesAllowed("user")    
    @Override
    @Blocking
    public Uni<proto.LocationTypeListResponse> listLocationTypes(Empty request) {
        return createUni(() -> LocationTypeController.listLocationTypes(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<Empty> deleteLocationType(proto.DeleteLocationTypeRequest request) {
        return createUni(() -> LocationTypeController.deleteLocationType(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<proto.LocationResponse> createLocation(proto.CreateLocationRequest request) {
        return createUni(() -> LocationController.createLocation(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<proto.LocationResponse> updateLocation(proto.UpdateLocationRequest request) {
        return createUni(() -> LocationController.updateLocation(request));
    }

    @RolesAllowed("user")        
    @Blocking
    @Override
    public Uni<proto.LocationListResponse> listLocations(proto.ListLocationsRequest request) {
        return createUni(() -> LocationController.listLocations(request));
    }        

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<Empty> deleteLocation(proto.DeleteLocationRequest request) {
        return createUni(() -> LocationController.deleteLocation(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<proto.NodeResponse> createNode(proto.CreateNodeRequest request) {
        return createUni(() -> NodeController.createNode(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<proto.NodeResponse> updateNode(proto.UpdateNodeRequest request) {
        return createUni(() -> NodeController.updateNode(request));
    }

    @RolesAllowed("user")    
    @Blocking
    @Override
    public Uni<proto.NodeListResponse> listNodes(proto.ListNodesRequest request) {
        return createUni(() -> NodeController.listNodes(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<Empty> deleteNode(proto.DeleteNodeRequest request) {
        return createUni(() -> NodeController.deleteNode(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<proto.DeviceResponse> createDevice(proto.CreateDeviceRequest request) {
        return createUni(() -> DeviceController.createDevice(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<proto.DeviceResponse> updateDevice(proto.UpdateDeviceRequest request) {
        return createUni(() -> DeviceController.updateDevice(request));
    }

    @RolesAllowed("user")    
    @Blocking
    @Override
    public Uni<proto.DeviceListResponse> listDevices(proto.ListDevicesRequest request) {
        return createUni(() -> DeviceController.listDevices(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<Empty> deleteDevice(proto.DeleteDeviceRequest request) {
        return createUni(() -> DeviceController.deleteDevice(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<proto.ChannelResponse> createChannel(proto.CreateChannelRequest request) {
        return createUni(() -> ChannelController.createChannel(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<proto.ChannelResponse> updateChannel(proto.UpdateChannelRequest request) {
        return createUni(() -> ChannelController.updateChannel(request));
    }

    @RolesAllowed("user")        
    @Blocking
    @Override
    public Uni<proto.ChannelListResponse> listChannels(proto.ListChannelsRequest request) {
        return createUni(() -> ChannelController.listChannels(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<Empty> deleteChannel(proto.DeleteChannelRequest request) {
        return createUni(() -> ChannelController.deleteChannel(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<ChannelAccessControlResponse> addChannelAccessControl(AddChannelAccessControlRequest request) {
        return createUni(() -> ChannelController.addChannelAccessControl(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<ChannelAlarmResponse> addChannelAlarm(AddChannelAlarmRequest request) {
        return createUni(() -> ChannelController.addChannelAlarm(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<ChannelTransformResponse> addChannelTransform(AddChannelTransformRequest request) {
        return createUni(() -> ChannelController.addChannelTransform(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<AlarmTypeResponse> createAlarmType(CreateAlarmTypeRequest request) {          
        return createUni(() -> AlarmTypeController.createAlarmType(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<RoleResponse> createRole(CreateRoleRequest request) {
        return createUni(() -> RoleController.createRole(request));
    }

    @RolesAllowed("admin")    
    @Blocking
    @Override
    public Uni<Empty> deleteAlarmType(DeleteAlarmTypeRequest request) {
        return createUni(() -> AlarmTypeController.deleteAlarmType(request));
    }

    @RolesAllowed("admin")
    @Blocking
    @Override
    public Uni<Empty> deleteChannelAccessControl(DeleteChannelAccessControlRequest request) {
        return createUni(() -> ChannelController.deleteChannelAccessControl(request));
    }

    @RolesAllowed("admin")    
    @Blocking
    @Override
    public Uni<Empty> deleteChannelAlarm(DeleteChannelAlarmRequest request) {
        return createUni(() -> ChannelController.deleteChannelAlarm(request));
    }

    @RolesAllowed("admin")    
    @Blocking
    @Override
    public Uni<Empty> deleteChannelTransform(DeleteChannelTransformRequest request) {
        return createUni(() -> ChannelController.deleteChannelTransform(request));
    }

    @RolesAllowed("admin")    
    @Blocking
    @Override
    public Uni<Empty> deleteRole(DeleteRoleRequest request) {
        return createUni(() -> RoleController.deleteRole(request));
    }

    @RolesAllowed("user")    
    @Blocking
    @Override
    public Uni<ChannelResponse> getChannel(GetChannelRequest request) {
       return createUni(() -> ChannelController.getChannel(request));
    }

    @RolesAllowed("user")        
    @Blocking
    @Override
    public Uni<DeviceResponse> getDevice(GetDeviceRequest request) {
        return createUni(() -> DeviceController.getDevice(request));
    }

    @RolesAllowed("user")        
    @Blocking
    @Override
    public Uni<LocationResponse> getLocation(GetLocationRequest request) {
        return createUni(() -> LocationController.getLocation(request));        
    }

    @RolesAllowed("user")        
    @Blocking
    @Override
    public Uni<NodeResponse> getNode(GetNodeRequest request) {
        return createUni(() -> NodeController.getNode(request));
    }

    @RolesAllowed("user")        
    @Blocking
    @Override
    public Uni<AlarmTypeListResponse> listAlarmTypes(Empty request) {
        return createUni(() -> AlarmTypeController.listAlarmTypes(request));
    }

    @RolesAllowed("user")        
    @Blocking
    @Override
    public Uni<RoleListResponse> listRoles(Empty request) {
        return createUni(() -> RoleController.listRoles(request));
    }

    @RolesAllowed("admin")    
    @Blocking
    @Override
    public Uni<AlarmTypeResponse> updateAlarmType(UpdateAlarmTypeRequest request) {
        return createUni(() -> AlarmTypeController.updateAlarmType(request));
    }

    @RolesAllowed("admin")    
    @Blocking
    @Override
    public Uni<ChannelAccessControlResponse> updateChannelAccessControl(UpdateChannelAccessControlRequest request) {
        return createUni(() -> ChannelController.updateChannelAccessControl(request));
    }

    @RolesAllowed("admin")    
    @Blocking
    @Override
    public Uni<ChannelAlarmResponse> updateChannelAlarm(UpdateChannelAlarmRequest request) {
        return createUni(() -> ChannelController.updateChannelAlarm(request));
    }

    @RolesAllowed("admin")    
    @Blocking
    @Override
    public Uni<ChannelTransformResponse> updateChannelTransform(UpdateChannelTransformRequest request) {
        return createUni(() -> ChannelController.updateChannelTransform(request));
    }

    @RolesAllowed("admin")    
    @Blocking
    @Override
    public Uni<RoleResponse> updateRole(UpdateRoleRequest request) {
        return createUni(() -> RoleController.updateRole(request));
    }

}
