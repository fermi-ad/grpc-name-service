package nameserver.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.validation.constraints.NotBlank;

@Entity
public class AlarmType extends PanacheEntity {
    @NotBlank
    @Column(unique = true)
    private String name;
    
    private String description;

    public AlarmType() {};

    public AlarmType(String name, String description) {
        this.name = name;
        this.description = description;
    }
    
    public String getName() {
        return name;
    }


    public AlarmType setName(String name) {
        this.name = name;
        return this;
    }


    public String getDescription() {
        return description;
    }


    public AlarmType setDescription(String description) {
        this.description = description;
        return this;
    }

}
