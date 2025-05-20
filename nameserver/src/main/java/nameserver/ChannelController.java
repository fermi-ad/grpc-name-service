package nameserver;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import com.google.protobuf.Empty;

import io.quarkus.hibernate.orm.panache.PanacheQuery;
import io.quarkus.panache.common.Parameters;
import jakarta.transaction.Transactional;
import nameserver.model.Channel;
import nameserver.model.ChannelAlarm;
import nameserver.model.ChannelTransform;
import nameserver.model.ChannelAccessControl;
import jakarta.enterprise.context.ApplicationScoped;

/*
 * Channel API implementation
 */
//make class application scoped to make the instance accessible by unit tests
@ApplicationScoped
public class ChannelController {
    
    public static Channel getChannelByName(String name) {        
        Channel channel = Channel.find("name", name).firstResult();
        if (channel == null) {
            throw new RuntimeException("Channel with name " + name + " not found");
        }
        return channel;
    }

    public static ChannelAlarm getChannelAlarm(String channel, String alarmType) {        
        ChannelAlarm channelAlarm = ChannelAlarm.find("channel.name = ?1 and alarmType.name = ?2", channel, alarmType).firstResult();
        if (channelAlarm == null) {
            throw new RuntimeException(alarmType + " alarm  for channel " + channel + " not found");
        }
        return channelAlarm;
    }
    
    public static ChannelTransform getChannelTransform(String channel, String transform) {        
        ChannelTransform channelTransform = ChannelTransform.find("channel.name = ?1 and name = ?2", channel, transform).firstResult();
        if (channelTransform == null) {
            throw new RuntimeException(transform + " transform for channel " + channel + " not found");
        }
        return channelTransform;
    }
    public static ChannelAccessControl getChannelAccessControl(String channel, String role) {        
        ChannelAccessControl channelAccessControl = ChannelAccessControl.find("channel.name = ?1 and role.name = ?2", channel, role).firstResult();
        if (channelAccessControl == null) {
            throw new RuntimeException(role + " access control for channel " + channel + " not found");
        }
        return channelAccessControl;
    }

    private static proto.ChannelAlarm.Builder createChannelAlarmBuilder(ChannelAlarm channelAlarm) {
        var channelAlarmBuilder = proto.ChannelAlarm.newBuilder();
        
        channelAlarmBuilder.setTriggerCondition(channelAlarm.getTriggerCondition())
                .setType(channelAlarm.getAlarmType().getName());
        
        return channelAlarmBuilder;
    }

    private static proto.ChannelTransform.Builder createChannelTransformBuilder(ChannelTransform channelTransform) {
        var channelTransformBuilder = proto.ChannelTransform.newBuilder();
        
        channelTransformBuilder.setName(channelTransform.getName())
                .setDescription(channelTransform.getDescription())
                .setTransform(channelTransform.getTransform());
        
        return channelTransformBuilder;
    }

    private static proto.ChannelAccessControl.Builder createChannelAccessControlBuilder(ChannelAccessControl channelAccessControl) {
        var channelAccessControlBuilder = proto.ChannelAccessControl.newBuilder();
        
        channelAccessControlBuilder.setRole(channelAccessControl.getRole().getName())
                .setRead(channelAccessControl.isRead())
                .setWrite(channelAccessControl.isWrite());
        
        return channelAccessControlBuilder;
    }

    private static ChannelAlarm createChannelAlarm(proto.ChannelAlarm reqChannelAlarm, Channel channel) {
        var channelAlarm= new ChannelAlarm();
        
        channelAlarm
            .setTriggerCondition(reqChannelAlarm.getTriggerCondition())
            .setAlarmType(AlarmTypeController.getAlarmTypeByName(reqChannelAlarm.getType()))
            .setChannel(channel);        

        return channelAlarm;
    }

    private static ChannelTransform createChannelTransform(proto.ChannelTransform reqChannelTransform, Channel channel) {
        var channelTransform = new ChannelTransform();
        
        channelTransform
            .setName(reqChannelTransform.getName())
            .setDescription(reqChannelTransform.getDescription())
            .setTransform(reqChannelTransform.getTransform())
            .setChannel(channel);
        
        return channelTransform;
    }

    private static ChannelAccessControl createChannelAccessControl(proto.ChannelAccessControl reqChannelAccessControl, Channel channel) {
        var channelAccessControl = new ChannelAccessControl();
        
        channelAccessControl
            .setRole(RoleController.getRoleByName(reqChannelAccessControl.getRole()))
            .setRead(reqChannelAccessControl.getRead())
            .setWrite(reqChannelAccessControl.getWrite())
            .setChannel(channel);
        
        return channelAccessControl;
    }

    private static ChannelAlarm updateChannelAlarm(proto.ChannelAlarm reqChannelAlarm, String channelName) {
        ChannelAlarm channelAlarm = getChannelAlarm(channelName, reqChannelAlarm.getType());

        if (reqChannelAlarm.getTriggerCondition() != "") {
            channelAlarm.setTriggerCondition(reqChannelAlarm.getTriggerCondition());
        }
        if (reqChannelAlarm.getType() != "") {
            var alarmType = AlarmTypeController.getAlarmTypeByName(reqChannelAlarm.getType());
            channelAlarm.setAlarmType(alarmType);
        }   
        return channelAlarm;
    }

    private static ChannelTransform updateChannelTransform(proto.ChannelTransform reqChannelTransform, String channelName) {
        ChannelTransform channelTransform = getChannelTransform(channelName, reqChannelTransform.getName());

        if (reqChannelTransform.getDescription() != "") {
            channelTransform.setDescription(reqChannelTransform.getDescription());
        }
        if (reqChannelTransform.getTransform() != "") {
            channelTransform.setTransform(reqChannelTransform.getTransform());
        }   
        return channelTransform;
    }

    private static ChannelAccessControl updateChannelAccessControl(proto.ChannelAccessControl reqChannelAccessControl, String channelName) {
        ChannelAccessControl channelAccessControl = getChannelAccessControl(channelName, reqChannelAccessControl.getRole());

        if(reqChannelAccessControl.hasRead()) {
            channelAccessControl.setRead(reqChannelAccessControl.getRead());
        }

        if(reqChannelAccessControl.hasWrite()) {
            channelAccessControl.setWrite(reqChannelAccessControl.getWrite());
        }

        return channelAccessControl;
    }

    private static proto.Channel.Builder createChannelBuilder(Channel channel) {
        var channelBuilder = proto.Channel.newBuilder();        
        channelBuilder.setName(channel.getName());
        if(channel.getDescription() != null) {
            channelBuilder.setDescription(channel.getDescription());
        }
        if(channel.getMetadata() != null) {
            channelBuilder.setMetadata(channel.getMetadata());
        }        
        channelBuilder.setDeviceName(channel.getDevice().getName());
        
        if(channel.getAlarms() != null) {            
            for( var alarm : channel.getAlarms()) {            
                var alarmBuilder = createChannelAlarmBuilder(alarm);
                channelBuilder.addAlarms(alarmBuilder);
            }
        }
        
        if(channel.getTransforms() != null) {
            for( var transform : channel.getTransforms()) {
                var transformBuilder = createChannelTransformBuilder(transform);
                channelBuilder.addTransforms(transformBuilder);
            }
        }

        if(channel.getAccessControls() != null) {
            for(var accessControl : channel.getAccessControls()) {
                var accessControlBuilder = createChannelAccessControlBuilder(accessControl);
                channelBuilder.addAccesscontrols(accessControlBuilder);
            }
        }

        return channelBuilder;
    } 

    @Transactional
    public static proto.ChannelResponse createChannel(proto.CreateChannelRequest request) {        
        var reqChannel = request.getChannel();
        
        var device = DeviceController.getDeviceByName(reqChannel.getDeviceName());

        Channel channel = new Channel(reqChannel.getName(),
                             reqChannel.getDescription(),
                             reqChannel.getMetadata(),
                             device);
        
        if (reqChannel.getAlarmsList() != null) {
            for(var alarm : reqChannel.getAlarmsList()) {
                var channelAlarm = createChannelAlarm(alarm, channel);
                channel.addAlarm(channelAlarm);
            }
        }
                
        if (reqChannel.getTransformsList() != null) {
            for(var transform : reqChannel.getTransformsList()) {            
                var channelTransform = createChannelTransform(transform, channel);
                channel.addTransform(channelTransform);            
            }
        }
        
        if(reqChannel.getAccesscontrolsList() != null) {
            for(var accessControl : reqChannel.getAccesscontrolsList()) {
                var channelAccessControl = createChannelAccessControl(accessControl, channel);              
                channel.addAccessControl(channelAccessControl);
            }
        }
        
        channel.persist();
        var resChannelBuilder = createChannelBuilder(channel);        
        return proto.ChannelResponse.newBuilder().setChannel(resChannelBuilder).build();
    }

    @Transactional
    public static proto.ChannelResponse updateChannel(proto.UpdateChannelRequest request) {
        var reqChannel = request.getChannel();
        Channel channel = getChannelByName(reqChannel.getName());        
        if(reqChannel.getDescription() != "") {
            channel.setDescription(reqChannel.getDescription());
        }
        if(reqChannel.getMetadata() != "") {
            channel.setMetadata(reqChannel.getMetadata());
        }
        if(reqChannel.getDeviceName() != "") {
            channel.setDevice(DeviceController.getDeviceByName(reqChannel.getDeviceName()));
        }        

        java.util.List<ChannelAlarm> alarms = new java.util.ArrayList<>();
        
        if(reqChannel.getAlarmsList() != null) {
            for(var alarm : reqChannel.getAlarmsList()) {
                var channelAlarm = updateChannelAlarm(alarm, channel.getName());            
            }
        }

        if(reqChannel.getTransformsList() != null) {
            for(var transform : reqChannel.getTransformsList()) {            
                var channelTransform = updateChannelTransform(transform, channel.getName());            
            }
        }
        if(reqChannel.getAccesscontrolsList() != null) {
            for(var accessControl : reqChannel.getAccesscontrolsList()) {
                var channelAccessControl = updateChannelAccessControl(accessControl, channel.getName());            
            }
        }
        channel.persist();
        var resChannelBuilder = createChannelBuilder(channel);        

        return proto.ChannelResponse.newBuilder().setChannel(resChannelBuilder).build();                
    }

    public static proto.ChannelListResponse listChannels(proto.ListChannelsRequest request) {
        var queryParams = new ArrayList<List<String>>();
        if (!request.getName().isEmpty()) {
            queryParams.add(Arrays.asList("name",request.getName()));
        }
        if (!request.getDeviceName().isEmpty()) {
            queryParams.add(Arrays.asList("device.name",request.getDeviceName()));
        }
        Parameters parameters = new Parameters();   
            
        var queryString = ControllerUtil.createFindParam(queryParams, parameters);        
        PanacheQuery<Channel> query;
        query = Channel.find(queryString, parameters);
        
        var pag = ControllerUtil.readPagination(request.getPagination());                
        List<Channel> channels = query.page(pag[0] - 1, pag[1]).list();        
        proto.ChannelListResponse.Builder response = proto.ChannelListResponse.newBuilder();
        for (Channel channel : channels) {
            var chanBuilder = createChannelBuilder(channel);            
            response.addChannels(chanBuilder.build());
        }
        response.setPagination(ControllerUtil.createPaginationResponse(pag, query.count()));
        return response.build();
    }

    @Transactional
    public static Empty deleteChannel(proto.DeleteChannelRequest request) {        
        Channel.delete("name", request.getName());
        return Empty.newBuilder().build();
    }     

    public static proto.ChannelResponse getChannel(proto.GetChannelRequest request) {
        Channel channel = getChannelByName(request.getName());
        var channelBuilder = createChannelBuilder(channel);
        return proto.ChannelResponse.newBuilder().setChannel(channelBuilder).build();
    }

    @Transactional
    public static proto.ChannelAlarmResponse addChannelAlarm(proto.AddChannelAlarmRequest request) {
        var reqAlarm = request.getAlarm();
        Channel channel = getChannelByName(request.getChannelName());
        var channelAlarm = createChannelAlarm(reqAlarm, channel);
        channelAlarm.persist();
        var alarmBuilder = createChannelAlarmBuilder(channelAlarm);
        return proto.ChannelAlarmResponse.newBuilder().setAlarm(alarmBuilder).build();
    }

    @Transactional
    public static proto.ChannelTransformResponse addChannelTransform(proto.AddChannelTransformRequest request) {
        var reqTransform = request.getTransform();
        Channel channel = getChannelByName(request.getChannelName());
        var channelTransform = createChannelTransform(reqTransform, channel);
        channelTransform.persist();
        var transformBuilder = createChannelTransformBuilder(channelTransform);
        return proto.ChannelTransformResponse.newBuilder().setTransform(transformBuilder).build();
    }

    @Transactional
    public static proto.ChannelAccessControlResponse addChannelAccessControl(proto.AddChannelAccessControlRequest request) {
        var reqAccessControl = request.getAccesscontrol();
        Channel channel = getChannelByName(request.getChannelName());
        var channelAccessControl = createChannelAccessControl(reqAccessControl, channel);
        channelAccessControl.persist();
        var accessControlBuilder = createChannelAccessControlBuilder(channelAccessControl);
        return proto.ChannelAccessControlResponse.newBuilder().setAccesscontrol(accessControlBuilder).build();
    }

    @Transactional
    public static Empty deleteChannelAlarm(proto.DeleteChannelAlarmRequest request) {     
        ChannelAlarm.delete("channel.name = ?1 and alarmType.name = ?2", request.getChannelName(), request.getAlarmType());
        return Empty.newBuilder().build();
    }

    @Transactional
    public static Empty deleteChannelTransform(proto.DeleteChannelTransformRequest request) {     
        ChannelTransform.delete("channel.name = ?1 and name = ?2", request.getChannelName(), request.getTransformName());
        return Empty.newBuilder().build();
    }

    @Transactional 
    public static Empty deleteChannelAccessControl(proto.DeleteChannelAccessControlRequest request) {     
        ChannelAccessControl.delete("channel.name = ?1 and role.name = ?2", request.getChannelName(), request.getRole());
        return Empty.newBuilder().build();
    }

    @Transactional
    public static proto.ChannelAlarmResponse updateChannelAlarm(proto.UpdateChannelAlarmRequest request) {        
        ChannelAlarm channelAlarm = updateChannelAlarm(request.getAlarm(), request.getChannelName());
        channelAlarm.persist();
        var alarmBuilder = createChannelAlarmBuilder(channelAlarm);
        return proto.ChannelAlarmResponse.newBuilder().setAlarm(alarmBuilder).build();
    }

    @Transactional
    public static proto.ChannelTransformResponse updateChannelTransform(proto.UpdateChannelTransformRequest request) {        
        ChannelTransform channelTransform = updateChannelTransform(request.getTransform(), request.getChannelName());
        channelTransform.persist();
        var transformBuilder = createChannelTransformBuilder(channelTransform);
        return proto.ChannelTransformResponse.newBuilder().setTransform(transformBuilder).build();
    }

    @Transactional
    public static proto.ChannelAccessControlResponse updateChannelAccessControl(proto.UpdateChannelAccessControlRequest request) {        
        ChannelAccessControl channelAccessControl = updateChannelAccessControl(request.getAccesscontrol(), request.getChannelName());
        channelAccessControl.persist();
        var accessControlBuilder = createChannelAccessControlBuilder(channelAccessControl);
        return proto.ChannelAccessControlResponse.newBuilder().setAccesscontrol(accessControlBuilder).build();
    }
}

