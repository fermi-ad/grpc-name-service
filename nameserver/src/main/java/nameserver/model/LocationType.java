package nameserver.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Table;
import jakarta.validation.constraints.NotBlank;

@Entity
public class LocationType extends PanacheEntity {

    public String getName() {
        return name;
    }

    public LocationType setName(String name) {
        this.name = name;
        return this;        
    }

    public String getDescription() {
        return description;
    }

    public LocationType setDescription(String description) {
        this.description = description;
        return this;        
    }

    @NotBlank
    @Column(unique = true)
    public String name;    
    public String description;


    public LocationType() {
        // Default constructor
    }

    public LocationType(String name, String description) {
        this.name = name;
        this.description = description;
    }
}
