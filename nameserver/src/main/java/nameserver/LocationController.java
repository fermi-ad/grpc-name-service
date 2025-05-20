package nameserver;
import java.util.ArrayList;
import java.util.List;

import io.quarkus.hibernate.orm.panache.PanacheQuery;
import io.quarkus.panache.common.Parameters;
import jakarta.transaction.Transactional;
import nameserver.model.Location;
import com.google.protobuf.Empty;
import java.util.Arrays;

/*Location API implementation */
public class LocationController  {
    
    public static Location getLocationByName(String name) {        
        Location location = Location.find("name", name).firstResult();
        if (location == null) {
            throw new RuntimeException("Location with name " + name + " not found");
        }
        return location;
    }
    
    static proto.Location.Builder createLocationBuilder(Location location) {
        var locationBuilder = proto.Location.newBuilder();
        
        locationBuilder.setName(location.getName())
                .setDescription(location.getDescription())
                .setLocationTypeName(location.getLocationType().getName());
        
        if (location.getParentLocation() != null) {
            locationBuilder.setParentLocationName(location.getParentLocation().getName());
        }
        
        return locationBuilder;
    }


    @Transactional
    public static proto.LocationResponse createLocation(proto.CreateLocationRequest request) {        
        var reqLoc = request.getLocation();        
        var locationType = LocationTypeController.getLocationTypeByName(reqLoc.getLocationTypeName());
        Location location = new Location(reqLoc.getName(), reqLoc.getDescription(), locationType);
        if (!reqLoc.getParentLocationName().isEmpty()) {
            Location parentLocation = getLocationByName(reqLoc.getParentLocationName());
            location.setParentLocation(parentLocation);
        }
        location.persist();
        var resLocBuilder = createLocationBuilder(location);
        return proto.LocationResponse.newBuilder().setLocation(resLocBuilder).build();
    }

    @Transactional
    public static proto.LocationResponse updateLocation(proto.UpdateLocationRequest request) {
        var reqLoc = request.getLocation();
        Location location = getLocationByName(reqLoc.getName()); 
        if (!reqLoc.getDescription().isEmpty()) {      
            location.setDescription(reqLoc.getDescription());
        }
        location.setLocationType(LocationTypeController.getLocationTypeByName(reqLoc.getLocationTypeName()));
        if (reqLoc.getParentLocationName().isEmpty()) {
            Location parentLocation = getLocationByName(reqLoc.getParentLocationName());
            location.setParentLocation(parentLocation);
        }

        location.persist();
        var resLocBuilder = createLocationBuilder(location);        

        return proto.LocationResponse.newBuilder().setLocation(resLocBuilder).build();                
    }

    public static proto.LocationListResponse listLocations(proto.ListLocationsRequest request) {          
        var queryParams = new ArrayList<List<String>>();
        if (!request.getName().isEmpty()) {
            queryParams.add(Arrays.asList("name",request.getName()));            
        }
        if (!request.getLocationTypeName().isEmpty()) {            
            queryParams.add(Arrays.asList("locationType.name",request.getLocationTypeName()));
        }
        if (!request.getParentLocationName().isEmpty()) {
            queryParams.add(Arrays.asList("parentLocation.name",request.getParentLocationName()));
        }
        Parameters parameters = new Parameters();   
            
        var queryString = ControllerUtil.createFindParam(queryParams, parameters);        
        PanacheQuery<Location> query;
        query = Location.find(queryString, parameters);
        
        var pag = ControllerUtil.readPagination(request.getPagination());                
        List<Location> locations = query.page(pag[0] - 1, pag[1]).list();        
        proto.LocationListResponse.Builder response = proto.LocationListResponse.newBuilder();
        for (Location location : locations) {
            var locBuilder = createLocationBuilder(location);            
            response.addLocations(locBuilder.build());
        }
        response.setPagination(ControllerUtil.createPaginationResponse(pag, query.count()));
        return response.build();
    }

    @Transactional
    public static Empty deleteLocation(proto.DeleteLocationRequest request) {        
        Location.delete("name", request.getName());
        return Empty.newBuilder().build();
    }     

    public static proto.LocationResponse getLocation(proto.GetLocationRequest request) {
        var location = getLocationByName(request.getName());
        var locBuilder = createLocationBuilder(location);
        return proto.LocationResponse.newBuilder().setLocation(locBuilder).build();
    }
}
