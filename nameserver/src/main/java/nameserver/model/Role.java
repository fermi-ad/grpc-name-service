package nameserver.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.validation.constraints.NotBlank;

@Entity
public class Role extends PanacheEntity {

    @NotBlank
    @Column(unique = true)
    private String name;    
    private String description;


    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Role() {
        // Default constructor
    }

    public Role(String name, String description) {
        this.name = name;
        this.description = description;
    }
}

