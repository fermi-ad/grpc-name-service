package nameserver;
import com.google.protobuf.Empty;

import jakarta.transaction.Transactional;
import nameserver.model.LocationType;
import java.util.List;
import io.grpc.StatusRuntimeException;
import io.grpc.Status;

/*Location Type API implementation */
public class LocationTypeController {
    public static LocationType getLocationTypeByName(String name) {        
        LocationType locationType = LocationType.find("name", name).firstResult();
        if (locationType == null) {
            throw new StatusRuntimeException(Status.NOT_FOUND.withDescription("LocationType with name " + name + " not found"));
        }
        return locationType;
    }

    static proto.LocationType.Builder createLocationTypeBuilder(LocationType locationType) {
        var locationTypeBuilder = proto.LocationType.newBuilder();
        
        locationTypeBuilder.setName(locationType.getName())
                .setDescription(locationType.getDescription());
        return locationTypeBuilder;
    }

    static void saveLocationType(LocationType locationType) {
        try {
            locationType.persist();
        } catch (Exception e) {
            throw new StatusRuntimeException(Status.UNKNOWN.withDescription("Error saving LocationType: " + e.getMessage()));
        }
    }
    @Transactional
    public static proto.LocationTypeResponse createLocationType(proto.CreateLocationTypeRequest request) {        
        var reqLocType = request.getLocationType();        
        var locationType = new LocationType(reqLocType.getName(), reqLocType.getDescription());
        saveLocationType(locationType);
        var resLocTypeBuilder = createLocationTypeBuilder(locationType);
        return proto.LocationTypeResponse.newBuilder().setLocationType(resLocTypeBuilder).build();
    }

    @Transactional
    public static proto.LocationTypeResponse updateLocationType(proto.UpdateLocationTypeRequest request) {
        var reqLocType = request.getLocationType();        
        var name = reqLocType.getName();
        var locationType = getLocationTypeByName(name);
        if (reqLocType.getDescription() != "") {
            locationType.setDescription(reqLocType.getDescription());
        }
        try {
            locationType.persist();
        } catch (Exception e) {
            throw new RuntimeException("Error updating LocationType: " + e.getMessage());
        }
        
        var resLocTypeBuilder = createLocationTypeBuilder(locationType);
        return proto.LocationTypeResponse.newBuilder().setLocationType(resLocTypeBuilder).build();        
    }

    public static proto.LocationTypeListResponse listLocationTypes(Empty request) {
        var responseBuilder = proto.LocationTypeListResponse.newBuilder();
        List<LocationType> locationTypes = LocationType.listAll();
        for (LocationType l : locationTypes) {
            var lbuilder = createLocationTypeBuilder(l);
            responseBuilder.addLocationTypes(lbuilder.build());
        }
        return responseBuilder.build();
    }

    @Transactional
    public static Empty deleteLocationType(proto.DeleteLocationTypeRequest request) {        
        LocationType.delete("name", request.getName());
        return Empty.newBuilder().build();
    }
}
