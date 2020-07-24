<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="2.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform" xmlns:graphml="http://graphml.graphdrawing.org/xmlns" xmlns:y="http://www.yworks.com/xml/graphml" xmlns="http://graphml.graphdrawing.org/xmlns" exclude-result-prefixes="#all">
	<xsl:output method="xml" version="1.0" encoding="UTF-8" indent="yes"/>
	
	<xsl:template match="@*|node()">
		<xsl:copy>
			<xsl:apply-templates select="@*|node()"/>
		</xsl:copy>
	</xsl:template>
	
	<xsl:template match="graphml:graphml">	
		<xsl:copy>
			<xsl:apply-templates select="@*"/>
			<key id="node-graphics" for="node" yfiles.type="nodegraphics"/>
			<key id="edge-graphics" for="edge" yfiles.type="edgegraphics"/>
			<xsl:apply-templates select="node()"/>
		</xsl:copy>
	</xsl:template>
	
	<xsl:template match="graphml:node">	
		<xsl:copy>
			<xsl:apply-templates select="@*|node()"/>
			<data key="node-graphics">
				<y:ShapeNode>
					<y:Geometry height="30.0" width="30.0" x="0.0" y="0.0"/>
					<y:Fill color="#00A82D" transparent="false"/>
					<y:BorderStyle hasColor="false" raised="false" type="line" width="1.0"/>
					<y:NodeLabel alignment="center" autoSizePolicy="content" borderDistance="1.0" fontFamily="Dialog" fontSize="10" fontStyle="plain" hasBackgroundColor="false" hasLineColor="false" height="16.2509765625" horizontalTextPosition="center" iconTextGap="4" modelName="sides" modelPosition="s" textColor="#000000" verticalTextPosition="bottom" visible="true" width="54.46875" x="-12.234375" xml:space="preserve" y="31.0"><xsl:value-of select="graphml:data[@key='node-label']"/></y:NodeLabel>
					<y:Shape type="ellipse"/>
				</y:ShapeNode>
			</data>
		</xsl:copy>
	</xsl:template>
</xsl:stylesheet>
