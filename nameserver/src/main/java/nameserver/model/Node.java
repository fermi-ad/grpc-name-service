package nameserver.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

@Entity
public class Node extends PanacheEntity {
    @NotBlank
    @Column(unique = true)
    private String hostname;
    
    @NotBlank
    @Column(name = "ip_address")
    private String ipAddress;
    
    private String description;

    @ManyToOne(optional=false)
    @JoinColumn(name = "location_id")
    private Location location;

    public String getHostname() {
        return hostname;
    }

    public void setHostname(String hostname) {
        this.hostname = hostname;
    }

    public String getIpAddress() {
        return ipAddress;
    }

    public void setIpAddress(String ipAddress) {
        this.ipAddress = ipAddress;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Location getLocation() {
        return location;
    }

    public void setLocation(Location location) {
        this.location = location;
    }

    public Node() {
        // Default constructor
    }

    public Node(String hostname, String ipAddress, Location location, String description) {
        this.hostname = hostname;
        this.ipAddress = ipAddress;
        this.location = location;
        this.description = description;
    }
}

