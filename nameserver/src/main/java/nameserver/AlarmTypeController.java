package nameserver;
import java.util.List;

import com.google.protobuf.Empty;

import jakarta.transaction.Transactional;
import nameserver.model.AlarmType;

public class AlarmTypeController {
    public static AlarmType getAlarmTypeByName(String name) {        
        AlarmType alarmType = AlarmType.find("name", name).firstResult();
        if (alarmType == null) {
            throw new RuntimeException("AlarmType with name " + name + " not found");
        }
        return alarmType;
    }

    static proto.AlarmType.Builder createAlarmTypeBuilder(AlarmType alarmType) {
        var alarmTypeBuilder = proto.AlarmType.newBuilder();
        
        alarmTypeBuilder.setName(alarmType.getName())
                .setDescription(alarmType.getDescription());
        return alarmTypeBuilder;
    }

    @Transactional
    public static proto.AlarmTypeResponse createAlarmType(proto.CreateAlarmTypeRequest request) {        
        var reqLocType = request.getAlarmType();        
        var alarmType = new AlarmType(reqLocType.getName(), reqLocType.getDescription());
        try {
            alarmType.persist();
        } catch (Exception e) {
            throw new RuntimeException("Error creating AlarmType: " + e.getMessage());
        }
        var resLocTypeBuilder = createAlarmTypeBuilder(alarmType);
        return proto.AlarmTypeResponse.newBuilder().setAlarmType(resLocTypeBuilder).build();
    }

    @Transactional
    public static proto.AlarmTypeResponse updateAlarmType(proto.UpdateAlarmTypeRequest request) {
        var reqLocType = request.getAlarmType();        
        var name = reqLocType.getName();
        var alarmType = getAlarmTypeByName(name);
        if (reqLocType.getDescription() != "") {
            alarmType.setDescription(reqLocType.getDescription());
        }
        try {
            alarmType.persist();
        } catch (Exception e) {
            throw new RuntimeException("Error updating AlarmType: " + e.getMessage());
        }
        
        var resLocTypeBuilder = createAlarmTypeBuilder(alarmType);
        return proto.AlarmTypeResponse.newBuilder().setAlarmType(resLocTypeBuilder).build();        
    }

    public static proto.AlarmTypeListResponse listAlarmTypes(Empty request) {
        var responseBuilder = proto.AlarmTypeListResponse.newBuilder();
        List<AlarmType> alarmTypes = AlarmType.listAll();
        for (AlarmType l : alarmTypes) {
            var lbuilder = createAlarmTypeBuilder(l);
            responseBuilder.addAlarmTypes(lbuilder.build());
        }
        return responseBuilder.build();
    }

    @Transactional
    public static Empty deleteAlarmType(proto.DeleteAlarmTypeRequest request) {        
       AlarmType.delete("name", request.getName());
       return Empty.newBuilder().build();
    }
}
