package nameserver.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

@Entity
public class Location extends PanacheEntity {
    @NotBlank
    @Column(unique = true)
    public String name;    
    public String description;

    @NotNull
    @ManyToOne
    @JoinColumn(name = "location_type_id")    
    public LocationType locationType;
    
    @ManyToOne(fetch=FetchType.LAZY)
    @JoinColumn(name = "parent_location_id")
    public Location parentLocation;

    public Location() {
        // Default constructor
    }

    public String getName() {
        return name;
    }

    public Location setName(String name) {
        this.name = name;
        return this;        
    }

    public String getDescription() {
        return description;
    }

    public Location setDescription(String description) {
        this.description = description;
        return this;        
    }

    public LocationType getLocationType() {
        return locationType;
    }

    public Location setLocationType(LocationType locationType) {
        this.locationType = locationType;
        return this;        
    }

    public Location getParentLocation() {
        return parentLocation;
    }

    public Location setParentLocation(Location parentLocation) {
        this.parentLocation = parentLocation;
        return this;        
    }

    public Location(String name, String description, LocationType locationType) {
        this.name = name;
        this.description = description;
        this.locationType = locationType;
        this.parentLocation = null;
    }
}

